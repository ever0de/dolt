package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	rand2 "crypto/rand"

	"github.com/dolthub/dolt/go/store/chunks"
	"github.com/dolthub/dolt/go/store/constants"
	"github.com/dolthub/dolt/go/store/hash"
	"github.com/dolthub/dolt/go/store/nbs"
	"github.com/dolthub/dolt/go/store/pool"
	"github.com/dolthub/dolt/go/store/prolly"
	"github.com/dolthub/dolt/go/store/prolly/message"
	"github.com/dolthub/dolt/go/store/prolly/tree"
	"github.com/dolthub/dolt/go/store/val"
	"golang.org/x/exp/rand"
)

var (
	r    = rand.New(rand.NewSource(1))
	keys [][]byte

	sharedPool = pool.NewBuffPool()

	keyDesc = val.NewTupleDescriptor(
		val.Type{Enc: val.ByteStringEnc, Nullable: false},
	)
	valDesc = val.NewTupleDescriptor(
		val.Type{Enc: val.ByteStringEnc, Nullable: true},
	)

	emptyNode tree.Node = func() tree.Node {
		s := message.NewProllyMapSerializer(val.TupleDesc{}, sharedPool)
		msg := s.Serialize(nil, nil, nil, 0)
		emptyNode, err := tree.NodeFromBytes(msg)
		if err != nil {
			panic(err)
		}

		return emptyNode
	}()
	emptyNodeHash = emptyNode.HashOf()

	emptyHash = hash.Hash{}
)

const (
	maxKeyLength    = 128
	maxValueLength  = 128
	operationSet    = 0
	operationDelete = 1
	totalKeys       = 50_000_000
)

func init() {
	keys = make([][]byte, totalKeys)
	for i := 0; i < totalKeys; i++ {
		seed := make([]byte, 16)
		rand2.Read(seed)
		hash := sha256.Sum256(seed)
		key := hash[:]

		if len(key) > maxKeyLength {
			key = key[:maxKeyLength]
		}

		keys[i] = key
	}
}

func main() {
	// cpu pprof
	f, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		panic(err)
	}

	defer pprof.StopCPUProfile()

	ctx := context.Background()
	path := "./tmp"
	os.MkdirAll(path, os.ModePerm)

	nbs, err := nbs.NewLocalJournalingStore(
		ctx,
		constants.FormatDoltString,
		path,
		nbs.NewUnlimitedMemQuotaProvider(),
	)
	if err != nil {
		panic(err)
	}

	rootHash, err := nbs.Root(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("root hash:", rootHash.String())

	chunkStore := chunks.ChunkStore(nbs)
	nodeStore := tree.NewNodeStore(chunkStore)

	versionRange := struct{ Start, End uint64 }{Start: 1, End: 800}

	m := prolly.NewMap(emptyNode, nodeStore, keyDesc, valDesc)

	setupCh := make(chan struct{})
	go func() {
		dbSetup(m, versionRange, 10_000)
		close(setupCh)
	}()

	sg := make(chan os.Signal, 1)
	signal.Notify(sg, os.Interrupt)
	select {
	case <-sg:
		return

	case <-setupCh:
	}
}

func dbSetup(m prolly.Map, versionRange struct{ Start, End uint64 }, changeSetCount int) {
	ctx := context.Background()
	var err error

	for i := versionRange.Start; i < versionRange.End; i++ {
		now := time.Now()
		mutable := m.Mutate()
		for j := 0; j < changeSetCount; j++ {
			key := keys[r.Intn(totalKeys)]

			err := mutable.Put(ctx, ProllyKey(key), val.EmptyTuple)
			if err != nil {
				panic(err)
			}
		}

		m, err = mutable.Map(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println("version", i, "took", time.Since(now).String())
	}
}

func ProllyKey(key []byte) val.Tuple {
	bld := val.NewTupleBuilder(keyDesc)
	bld.PutByteString(0, key)

	return bld.Build(sharedPool)
}
