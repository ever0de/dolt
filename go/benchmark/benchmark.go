package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime/pprof"
	"sort"
	"time"

	"github.com/dolthub/dolt/go/store/chunks"
	"github.com/dolthub/dolt/go/store/constants"
	"github.com/dolthub/dolt/go/store/nbs"
	"github.com/dolthub/dolt/go/store/pool"
	"github.com/dolthub/dolt/go/store/prolly"
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
)

const (
	maxKeyLength    = 128
	operationSet    = 0
	operationDelete = 1
	totalKeys       = 50_000_000
	// totalKeys = 1_000_000
)

func init() {
	keys = make([][]byte, totalKeys)

	retryCnt := 0
	for i := 0; i < totalKeys; i++ {
		key := make([]byte, maxKeyLength)
		// last 8 bytes to bigendian uint64
		key[maxKeyLength-8] = byte(i >> 56)
		key[maxKeyLength-7] = byte(i >> 48)
		key[maxKeyLength-6] = byte(i >> 40)
		key[maxKeyLength-5] = byte(i >> 32)
		key[maxKeyLength-4] = byte(i >> 24)
		key[maxKeyLength-3] = byte(i >> 16)
		key[maxKeyLength-2] = byte(i >> 8)
		key[maxKeyLength-1] = byte(i)

		keys[i] = key
	}

	fmt.Println("retry count:", retryCnt)
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

	nbs, err := nbs.NewLocalStore(
		// nbs, err := nbs.NewLocalJournalingStore(
		ctx,
		constants.FormatDoltString,
		path,
		128<<20,
		nbs.NewUnlimitedMemQuotaProvider(),
	)
	if err != nil {
		panic(err)
	}

	rootHash, err := nbs.Root(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("(nbs)root hash:", rootHash.String())

	chunkStore := chunks.ChunkStore(nbs)
	nodeStore := tree.NewNodeStore(chunkStore)

	versionRange := struct{ Start, End uint64 }{Start: 1, End: 100}

	m := dbSetup(nodeStore)
	cnt, _ := m.Count()
	fmt.Println("prolly map count:", cnt)
	fmt.Println("prolly map height:", m.Height())

	setupCh := make(chan struct{})
	go func() {
		set(m, versionRange, 2_000)
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

func dbSetup(ns tree.NodeStore) prolly.Map {
	kb := val.NewTupleBuilder(keyDesc)
	// sorted keys
	tups := make([]val.Tuple, totalKeys*2)
	for i := 0; i < totalKeys*2; i += 2 {
		kb.PutByteString(0, keys[i/2])
		tups[i] = kb.Build(sharedPool)
		tups[i+1] = val.EmptyTuple
	}

	m, err := prolly.NewMapFromTuples(context.TODO(), ns, keyDesc, valDesc, tups...)
	if err != nil {
		panic(err)
	}

	fmt.Println("created prolly map", m.HashOf().String())

	return m
}

func set(m prolly.Map, versionRange struct{ Start, End uint64 }, changeSetCount int) {
	ctx := context.Background()
	var err error

	// mem pprof
	memf, err := os.Create("mem.pprof")
	if err != nil {
		panic(err)
	}
	defer memf.Close()

	sinces := make([]time.Duration, 0)
	for i := versionRange.Start; i < versionRange.End; i++ {
		now := time.Now()
		mutable := m.Mutate()
		for j := 0; j < changeSetCount; j++ {
			idx := r.Intn(totalKeys)
			key := keys[idx]

			value := make([]byte, 1024)
			rand.Read(value)

			operation := r.Intn(2)
			if operation == 0 {
				err = mutable.Put(ctx, ProllyKey(key), ProllyKey(value))
			} else {
				err = mutable.Delete(ctx, ProllyKey(key))
			}
			if err != nil {
				panic(err)
			}
		}

		m, err = mutable.Map(ctx)
		if err != nil {
			panic(err)
		}

		since := time.Since(now)
		fmt.Println("version", i, "took", since)
		cnt, _ := m.Count()
		height := m.Height()
		fmt.Println("	prolly map count:", cnt, "height:", height)
		sinces = append(sinces, since)
	}

	err = pprof.WriteHeapProfile(memf)
	if err != nil {
		panic(err)
	}

	// avg, max, min, p95, p99, std dev
	avg, max, min, p95, p99, stdDev := stats(sinces)
	fmt.Println("avg:", avg)
	fmt.Println("max:", max)
	fmt.Println("min:", min)
	fmt.Println("95th percentile:", p95)
	fmt.Println("99th percentile:", p99)
	fmt.Println("std dev:", stdDev)
}

func stats(durations []time.Duration) (avg, max, min, p95, p99, stdDev time.Duration) {
	if len(durations) == 0 {
		return 0, 0, 0, 0, 0, 0
	}

	var sum time.Duration
	min = durations[0]
	max = durations[0]
	for _, duration := range durations {
		sum += duration
		if duration < min {
			min = duration
		}
		if duration > max {
			max = duration
		}
	}

	avg = sum / time.Duration(len(durations))

	sortedDurations := make([]time.Duration, len(durations))
	copy(sortedDurations, durations)
	sort.Slice(sortedDurations, func(i, j int) bool {
		return sortedDurations[i] < sortedDurations[j]
	})

	p95 = percentile(sortedDurations, 95)
	p99 = percentile(sortedDurations, 99)

	varianceSum := float64(0)
	for _, duration := range durations {
		diff := float64(duration - avg)
		varianceSum += diff * diff
	}

	variance := varianceSum / float64(len(durations))
	stdDev = time.Duration(math.Sqrt(variance))

	return avg, max, min, p95, p99, stdDev
}

func percentile(durations []time.Duration, percentile int) time.Duration {
	index := (percentile * len(durations) / 100) - 1
	if index < 0 {
		index = 0
	}
	return durations[index]
}
func ProllyKey(key []byte) val.Tuple {
	bld := val.NewTupleBuilder(keyDesc)
	bld.PutByteString(0, key)

	return bld.Build(sharedPool)
}
