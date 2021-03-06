// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hunts.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ForemanCheckin struct {
	LastHuntTimestamp     uint64   `protobuf:"varint,1,opt,name=last_hunt_timestamp,json=lastHuntTimestamp,proto3" json:"last_hunt_timestamp,omitempty"`
	LastEventTableVersion uint64   `protobuf:"varint,2,opt,name=last_event_table_version,json=lastEventTableVersion,proto3" json:"last_event_table_version,omitempty"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_unrecognized      []byte   `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *ForemanCheckin) Reset()         { *m = ForemanCheckin{} }
func (m *ForemanCheckin) String() string { return proto.CompactTextString(m) }
func (*ForemanCheckin) ProtoMessage()    {}
func (*ForemanCheckin) Descriptor() ([]byte, []int) {
	return fileDescriptor_0cb9fe450e8d0e17, []int{0}
}

func (m *ForemanCheckin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForemanCheckin.Unmarshal(m, b)
}
func (m *ForemanCheckin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForemanCheckin.Marshal(b, m, deterministic)
}
func (m *ForemanCheckin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForemanCheckin.Merge(m, src)
}
func (m *ForemanCheckin) XXX_Size() int {
	return xxx_messageInfo_ForemanCheckin.Size(m)
}
func (m *ForemanCheckin) XXX_DiscardUnknown() {
	xxx_messageInfo_ForemanCheckin.DiscardUnknown(m)
}

var xxx_messageInfo_ForemanCheckin proto.InternalMessageInfo

func (m *ForemanCheckin) GetLastHuntTimestamp() uint64 {
	if m != nil {
		return m.LastHuntTimestamp
	}
	return 0
}

func (m *ForemanCheckin) GetLastEventTableVersion() uint64 {
	if m != nil {
		return m.LastEventTableVersion
	}
	return 0
}

func init() {
	proto.RegisterType((*ForemanCheckin)(nil), "proto.ForemanCheckin")
}

func init() {
	proto.RegisterFile("hunts.proto", fileDescriptor_0cb9fe450e8d0e17)
}

var fileDescriptor_0cb9fe450e8d0e17 = []byte{
	// 186 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0xce, 0xb1, 0x8b, 0xc2, 0x30,
	0x14, 0x06, 0x70, 0x7a, 0xdc, 0xdd, 0x90, 0x83, 0x83, 0xeb, 0x21, 0x74, 0x14, 0x27, 0xa7, 0x64,
	0x28, 0xe2, 0xae, 0x28, 0xce, 0x52, 0x1c, 0x5c, 0x4a, 0x1a, 0x1f, 0x6d, 0x30, 0xc9, 0x2b, 0xc9,
	0x6b, 0xab, 0xff, 0xbd, 0x24, 0xc5, 0xe9, 0x83, 0xf7, 0xfb, 0x3e, 0x78, 0xec, 0xa7, 0x1b, 0x1c,
	0x05, 0xde, 0x7b, 0x24, 0xcc, 0xbf, 0x52, 0xac, 0x9e, 0xec, 0xf7, 0x88, 0x1e, 0xac, 0x74, 0xfb,
	0x0e, 0xd4, 0x5d, 0xbb, 0x9c, 0xb3, 0x7f, 0x23, 0x03, 0xd5, 0xb1, 0x5c, 0x93, 0xb6, 0x10, 0x48,
	0xda, 0xbe, 0xc8, 0x96, 0xd9, 0xfa, 0xf3, 0xfc, 0x17, 0xe9, 0x34, 0x38, 0xaa, 0xde, 0x90, 0x6f,
	0x59, 0x91, 0xfa, 0x30, 0x42, 0x1c, 0xc8, 0xc6, 0x40, 0x3d, 0x82, 0x0f, 0x1a, 0x5d, 0xf1, 0x91,
	0x46, 0x8b, 0xe8, 0x87, 0xc8, 0x55, 0xd4, 0xcb, 0x8c, 0xbb, 0xcd, 0xb5, 0x9c, 0xa6, 0x89, 0x8f,
	0x60, 0x50, 0xe9, 0x1b, 0x3c, 0xb8, 0x42, 0x2b, 0x5a, 0x34, 0xd2, 0xb5, 0x62, 0x3e, 0x7a, 0xd9,
	0x13, 0x7a, 0x21, 0x15, 0x69, 0x74, 0x41, 0xa4, 0x8f, 0x9b, 0xef, 0x14, 0xe5, 0x2b, 0x00, 0x00,
	0xff, 0xff, 0x64, 0x04, 0xab, 0x79, 0xce, 0x00, 0x00, 0x00,
}
