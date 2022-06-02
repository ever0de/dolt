// Copyright 2022 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package serial

import (
	"strconv"

	flatbuffers "github.com/google/flatbuffers/go"
)

type ItemType byte

const (
	ItemTypeUnknown          ItemType = 0
	ItemTypeTupleFormatAlpha ItemType = 1
)

var EnumNamesItemType = map[ItemType]string{
	ItemTypeUnknown:          "Unknown",
	ItemTypeTupleFormatAlpha: "TupleFormatAlpha",
}

var EnumValuesItemType = map[string]ItemType{
	"Unknown":          ItemTypeUnknown,
	"TupleFormatAlpha": ItemTypeTupleFormatAlpha,
}

func (v ItemType) String() string {
	if s, ok := EnumNamesItemType[v]; ok {
		return s
	}
	return "ItemType(" + strconv.FormatInt(int64(v), 10) + ")"
}

type ProllyTreeNode struct {
	_tab flatbuffers.Table
}

func GetRootAsProllyTreeNode(buf []byte, offset flatbuffers.UOffsetT) *ProllyTreeNode {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ProllyTreeNode{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsProllyTreeNode(buf []byte, offset flatbuffers.UOffsetT) *ProllyTreeNode {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ProllyTreeNode{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *ProllyTreeNode) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ProllyTreeNode) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ProllyTreeNode) KeyItems(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *ProllyTreeNode) KeyItemsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ProllyTreeNode) KeyItemsBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ProllyTreeNode) MutateKeyItems(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *ProllyTreeNode) KeyOffsets(j int) uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint16(a + flatbuffers.UOffsetT(j*2))
	}
	return 0
}

func (rcv *ProllyTreeNode) KeyOffsetsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ProllyTreeNode) MutateKeyOffsets(j int, n uint16) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint16(a+flatbuffers.UOffsetT(j*2), n)
	}
	return false
}

func (rcv *ProllyTreeNode) KeyType() ItemType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return ItemType(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *ProllyTreeNode) MutateKeyType(n ItemType) bool {
	return rcv._tab.MutateByteSlot(8, byte(n))
}

func (rcv *ProllyTreeNode) ValueItems(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *ProllyTreeNode) ValueItemsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ProllyTreeNode) ValueItemsBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ProllyTreeNode) MutateValueItems(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *ProllyTreeNode) ValueOffsets(j int) uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint16(a + flatbuffers.UOffsetT(j*2))
	}
	return 0
}

func (rcv *ProllyTreeNode) ValueOffsetsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ProllyTreeNode) MutateValueOffsets(j int, n uint16) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint16(a+flatbuffers.UOffsetT(j*2), n)
	}
	return false
}

func (rcv *ProllyTreeNode) ValueType() ItemType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return ItemType(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *ProllyTreeNode) MutateValueType(n ItemType) bool {
	return rcv._tab.MutateByteSlot(14, byte(n))
}

func (rcv *ProllyTreeNode) ValueAddressOffsets(j int) uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint16(a + flatbuffers.UOffsetT(j*2))
	}
	return 0
}

func (rcv *ProllyTreeNode) ValueAddressOffsetsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ProllyTreeNode) MutateValueAddressOffsets(j int, n uint16) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint16(a+flatbuffers.UOffsetT(j*2), n)
	}
	return false
}

func (rcv *ProllyTreeNode) AddressArray(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *ProllyTreeNode) AddressArrayLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ProllyTreeNode) AddressArrayBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ProllyTreeNode) MutateAddressArray(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *ProllyTreeNode) SubtreeCounts(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *ProllyTreeNode) SubtreeCountsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ProllyTreeNode) SubtreeCountsBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ProllyTreeNode) MutateSubtreeCounts(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *ProllyTreeNode) TreeCount() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ProllyTreeNode) MutateTreeCount(n uint64) bool {
	return rcv._tab.MutateUint64Slot(22, n)
}

func (rcv *ProllyTreeNode) TreeLevel() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ProllyTreeNode) MutateTreeLevel(n byte) bool {
	return rcv._tab.MutateByteSlot(24, n)
}

func ProllyTreeNodeStart(builder *flatbuffers.Builder) {
	builder.StartObject(11)
}
func ProllyTreeNodeAddKeyItems(builder *flatbuffers.Builder, keyItems flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(keyItems), 0)
}
func ProllyTreeNodeStartKeyItemsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func ProllyTreeNodeAddKeyOffsets(builder *flatbuffers.Builder, keyOffsets flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(keyOffsets), 0)
}
func ProllyTreeNodeStartKeyOffsetsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(2, numElems, 2)
}
func ProllyTreeNodeAddKeyType(builder *flatbuffers.Builder, keyType ItemType) {
	builder.PrependByteSlot(2, byte(keyType), 0)
}
func ProllyTreeNodeAddValueItems(builder *flatbuffers.Builder, valueItems flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(valueItems), 0)
}
func ProllyTreeNodeStartValueItemsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func ProllyTreeNodeAddValueOffsets(builder *flatbuffers.Builder, valueOffsets flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(valueOffsets), 0)
}
func ProllyTreeNodeStartValueOffsetsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(2, numElems, 2)
}
func ProllyTreeNodeAddValueType(builder *flatbuffers.Builder, valueType ItemType) {
	builder.PrependByteSlot(5, byte(valueType), 0)
}
func ProllyTreeNodeAddValueAddressOffsets(builder *flatbuffers.Builder, valueAddressOffsets flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(valueAddressOffsets), 0)
}
func ProllyTreeNodeStartValueAddressOffsetsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(2, numElems, 2)
}
func ProllyTreeNodeAddAddressArray(builder *flatbuffers.Builder, addressArray flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(addressArray), 0)
}
func ProllyTreeNodeStartAddressArrayVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func ProllyTreeNodeAddSubtreeCounts(builder *flatbuffers.Builder, subtreeCounts flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(subtreeCounts), 0)
}
func ProllyTreeNodeStartSubtreeCountsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func ProllyTreeNodeAddTreeCount(builder *flatbuffers.Builder, treeCount uint64) {
	builder.PrependUint64Slot(9, treeCount, 0)
}
func ProllyTreeNodeAddTreeLevel(builder *flatbuffers.Builder, treeLevel byte) {
	builder.PrependByteSlot(10, treeLevel, 0)
}
func ProllyTreeNodeEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type CommitClosure struct {
	_tab flatbuffers.Table
}

func GetRootAsCommitClosure(buf []byte, offset flatbuffers.UOffsetT) *CommitClosure {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &CommitClosure{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsCommitClosure(buf []byte, offset flatbuffers.UOffsetT) *CommitClosure {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &CommitClosure{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *CommitClosure) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *CommitClosure) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *CommitClosure) RefArray(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *CommitClosure) RefArrayLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *CommitClosure) RefArrayBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *CommitClosure) MutateRefArray(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *CommitClosure) TreeCount() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *CommitClosure) MutateTreeCount(n uint64) bool {
	return rcv._tab.MutateUint64Slot(6, n)
}

func (rcv *CommitClosure) TreeLevel() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *CommitClosure) MutateTreeLevel(n byte) bool {
	return rcv._tab.MutateByteSlot(8, n)
}

func CommitClosureStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func CommitClosureAddRefArray(builder *flatbuffers.Builder, refArray flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(refArray), 0)
}
func CommitClosureStartRefArrayVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func CommitClosureAddTreeCount(builder *flatbuffers.Builder, treeCount uint64) {
	builder.PrependUint64Slot(1, treeCount, 0)
}
func CommitClosureAddTreeLevel(builder *flatbuffers.Builder, treeLevel byte) {
	builder.PrependByteSlot(2, treeLevel, 0)
}
func CommitClosureEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
