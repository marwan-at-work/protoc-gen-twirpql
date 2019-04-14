// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

package e2e

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

type TrafficLight int32

const (
	TrafficLight_RED    TrafficLight = 0
	TrafficLight_YELLOW TrafficLight = 1
	TrafficLight_GREEN  TrafficLight = 2
)

var TrafficLight_name = map[int32]string{
	0: "RED",
	1: "YELLOW",
	2: "GREEN",
}

var TrafficLight_value = map[string]int32{
	"RED":    0,
	"YELLOW": 1,
	"GREEN":  2,
}

func (x TrafficLight) String() string {
	return proto.EnumName(TrafficLight_name, int32(x))
}

func (TrafficLight) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}

type HelloReq struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloReq) Reset()         { *m = HelloReq{} }
func (m *HelloReq) String() string { return proto.CompactTextString(m) }
func (*HelloReq) ProtoMessage()    {}
func (*HelloReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}

func (m *HelloReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloReq.Unmarshal(m, b)
}
func (m *HelloReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloReq.Marshal(b, m, deterministic)
}
func (m *HelloReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloReq.Merge(m, src)
}
func (m *HelloReq) XXX_Size() int {
	return xxx_messageInfo_HelloReq.Size(m)
}
func (m *HelloReq) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloReq.DiscardUnknown(m)
}

var xxx_messageInfo_HelloReq proto.InternalMessageInfo

func (m *HelloReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type HelloResp struct {
	Text                 string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloResp) Reset()         { *m = HelloResp{} }
func (m *HelloResp) String() string { return proto.CompactTextString(m) }
func (*HelloResp) ProtoMessage()    {}
func (*HelloResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{1}
}

func (m *HelloResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloResp.Unmarshal(m, b)
}
func (m *HelloResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloResp.Marshal(b, m, deterministic)
}
func (m *HelloResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloResp.Merge(m, src)
}
func (m *HelloResp) XXX_Size() int {
	return xxx_messageInfo_HelloResp.Size(m)
}
func (m *HelloResp) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloResp.DiscardUnknown(m)
}

var xxx_messageInfo_HelloResp proto.InternalMessageInfo

func (m *HelloResp) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type TrafficJamReq struct {
	Color                TrafficLight `protobuf:"varint,1,opt,name=color,proto3,enum=e2e.TrafficLight" json:"color,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *TrafficJamReq) Reset()         { *m = TrafficJamReq{} }
func (m *TrafficJamReq) String() string { return proto.CompactTextString(m) }
func (*TrafficJamReq) ProtoMessage()    {}
func (*TrafficJamReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{2}
}

func (m *TrafficJamReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TrafficJamReq.Unmarshal(m, b)
}
func (m *TrafficJamReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TrafficJamReq.Marshal(b, m, deterministic)
}
func (m *TrafficJamReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TrafficJamReq.Merge(m, src)
}
func (m *TrafficJamReq) XXX_Size() int {
	return xxx_messageInfo_TrafficJamReq.Size(m)
}
func (m *TrafficJamReq) XXX_DiscardUnknown() {
	xxx_messageInfo_TrafficJamReq.DiscardUnknown(m)
}

var xxx_messageInfo_TrafficJamReq proto.InternalMessageInfo

func (m *TrafficJamReq) GetColor() TrafficLight {
	if m != nil {
		return m.Color
	}
	return TrafficLight_RED
}

type TrafficJamResp struct {
	Next                 TrafficLight `protobuf:"varint,1,opt,name=next,proto3,enum=e2e.TrafficLight" json:"next,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *TrafficJamResp) Reset()         { *m = TrafficJamResp{} }
func (m *TrafficJamResp) String() string { return proto.CompactTextString(m) }
func (*TrafficJamResp) ProtoMessage()    {}
func (*TrafficJamResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{3}
}

func (m *TrafficJamResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TrafficJamResp.Unmarshal(m, b)
}
func (m *TrafficJamResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TrafficJamResp.Marshal(b, m, deterministic)
}
func (m *TrafficJamResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TrafficJamResp.Merge(m, src)
}
func (m *TrafficJamResp) XXX_Size() int {
	return xxx_messageInfo_TrafficJamResp.Size(m)
}
func (m *TrafficJamResp) XXX_DiscardUnknown() {
	xxx_messageInfo_TrafficJamResp.DiscardUnknown(m)
}

var xxx_messageInfo_TrafficJamResp proto.InternalMessageInfo

func (m *TrafficJamResp) GetNext() TrafficLight {
	if m != nil {
		return m.Next
	}
	return TrafficLight_RED
}

func init() {
	proto.RegisterEnum("e2e.TrafficLight", TrafficLight_name, TrafficLight_value)
	proto.RegisterType((*HelloReq)(nil), "e2e.HelloReq")
	proto.RegisterType((*HelloResp)(nil), "e2e.HelloResp")
	proto.RegisterType((*TrafficJamReq)(nil), "e2e.TrafficJamReq")
	proto.RegisterType((*TrafficJamResp)(nil), "e2e.TrafficJamResp")
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_a0b84a42fa06f626) }

var fileDescriptor_a0b84a42fa06f626 = []byte{
	// 247 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0x8d, 0x6d, 0x5a, 0xf3, 0x30, 0x21, 0x8e, 0x17, 0xe9, 0x41, 0x25, 0xe0, 0x1f, 0x3c,
	0xe4, 0x10, 0x11, 0x3d, 0x8b, 0x41, 0x91, 0xa0, 0xb0, 0x0a, 0xa2, 0xb7, 0x18, 0xa6, 0x36, 0x90,
	0x76, 0xd7, 0x6c, 0x10, 0x3f, 0xbe, 0xec, 0xa6, 0xc1, 0x05, 0xbd, 0xcd, 0xcc, 0x7b, 0xf3, 0x98,
	0xdf, 0x20, 0xd4, 0xdc, 0x7e, 0xd5, 0x15, 0xa7, 0xaa, 0x95, 0x9d, 0xa4, 0x11, 0x67, 0x9c, 0xec,
	0x63, 0xeb, 0x8e, 0x9b, 0x46, 0x0a, 0xfe, 0x24, 0xc2, 0x78, 0x55, 0x2e, 0x79, 0xcf, 0x3b, 0xf4,
	0x4e, 0x03, 0x61, 0xeb, 0xe4, 0x00, 0xc1, 0x5a, 0xd7, 0xca, 0x18, 0x3a, 0xfe, 0xee, 0x06, 0x83,
	0xa9, 0x93, 0x2b, 0x84, 0xcf, 0x6d, 0x39, 0x9f, 0xd7, 0xd5, 0x7d, 0xb9, 0x34, 0x29, 0x27, 0xf0,
	0x2b, 0xd9, 0xc8, 0xd6, 0xba, 0xa2, 0x6c, 0x27, 0xe5, 0x8c, 0xd3, 0xb5, 0xa5, 0xa8, 0x3f, 0x16,
	0x9d, 0xe8, 0xf5, 0xe4, 0x12, 0x91, 0xbb, 0xa9, 0x15, 0x1d, 0x61, 0xbc, 0x1a, 0xf2, 0xff, 0xdd,
	0xb4, 0xf2, 0x59, 0x8a, 0x6d, 0x77, 0x4a, 0x53, 0x8c, 0x44, 0x7e, 0x13, 0x6f, 0x10, 0x30, 0x79,
	0xcd, 0x8b, 0xe2, 0xf1, 0x25, 0xf6, 0x28, 0x80, 0x7f, 0x2b, 0xf2, 0xfc, 0x21, 0xde, 0xcc, 0x16,
	0x98, 0x3e, 0xf5, 0xe4, 0x74, 0x0c, 0xdf, 0xe2, 0x50, 0x68, 0xc3, 0x07, 0xf4, 0x59, 0xe4, 0xb6,
	0x5a, 0xd1, 0x05, 0xf0, 0x7b, 0x1b, 0x91, 0x7b, 0x49, 0x8f, 0x39, 0xdb, 0xfd, 0x33, 0xd3, 0xea,
	0xda, 0x7f, 0x33, 0x4f, 0x7d, 0x9f, 0xd8, 0x07, 0x9f, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0xe9,
	0x09, 0x52, 0x04, 0x71, 0x01, 0x00, 0x00,
}
