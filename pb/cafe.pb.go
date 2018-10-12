// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cafe.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CafeChallenge struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CafeChallenge) Reset()         { *m = CafeChallenge{} }
func (m *CafeChallenge) String() string { return proto.CompactTextString(m) }
func (*CafeChallenge) ProtoMessage()    {}
func (*CafeChallenge) Descriptor() ([]byte, []int) {
	return fileDescriptor_cafe_93e3fd513c792df6, []int{0}
}
func (m *CafeChallenge) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CafeChallenge.Unmarshal(m, b)
}
func (m *CafeChallenge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CafeChallenge.Marshal(b, m, deterministic)
}
func (dst *CafeChallenge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CafeChallenge.Merge(dst, src)
}
func (m *CafeChallenge) XXX_Size() int {
	return xxx_messageInfo_CafeChallenge.Size(m)
}
func (m *CafeChallenge) XXX_DiscardUnknown() {
	xxx_messageInfo_CafeChallenge.DiscardUnknown(m)
}

var xxx_messageInfo_CafeChallenge proto.InternalMessageInfo

func (m *CafeChallenge) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type CafeNonce struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CafeNonce) Reset()         { *m = CafeNonce{} }
func (m *CafeNonce) String() string { return proto.CompactTextString(m) }
func (*CafeNonce) ProtoMessage()    {}
func (*CafeNonce) Descriptor() ([]byte, []int) {
	return fileDescriptor_cafe_93e3fd513c792df6, []int{1}
}
func (m *CafeNonce) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CafeNonce.Unmarshal(m, b)
}
func (m *CafeNonce) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CafeNonce.Marshal(b, m, deterministic)
}
func (dst *CafeNonce) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CafeNonce.Merge(dst, src)
}
func (m *CafeNonce) XXX_Size() int {
	return xxx_messageInfo_CafeNonce.Size(m)
}
func (m *CafeNonce) XXX_DiscardUnknown() {
	xxx_messageInfo_CafeNonce.DiscardUnknown(m)
}

var xxx_messageInfo_CafeNonce proto.InternalMessageInfo

func (m *CafeNonce) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type CafeRegistration struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Nonce                string   `protobuf:"bytes,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Sig                  []byte   `protobuf:"bytes,4,opt,name=sig,proto3" json:"sig,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CafeRegistration) Reset()         { *m = CafeRegistration{} }
func (m *CafeRegistration) String() string { return proto.CompactTextString(m) }
func (*CafeRegistration) ProtoMessage()    {}
func (*CafeRegistration) Descriptor() ([]byte, []int) {
	return fileDescriptor_cafe_93e3fd513c792df6, []int{2}
}
func (m *CafeRegistration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CafeRegistration.Unmarshal(m, b)
}
func (m *CafeRegistration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CafeRegistration.Marshal(b, m, deterministic)
}
func (dst *CafeRegistration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CafeRegistration.Merge(dst, src)
}
func (m *CafeRegistration) XXX_Size() int {
	return xxx_messageInfo_CafeRegistration.Size(m)
}
func (m *CafeRegistration) XXX_DiscardUnknown() {
	xxx_messageInfo_CafeRegistration.DiscardUnknown(m)
}

var xxx_messageInfo_CafeRegistration proto.InternalMessageInfo

func (m *CafeRegistration) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *CafeRegistration) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *CafeRegistration) GetNonce() string {
	if m != nil {
		return m.Nonce
	}
	return ""
}

func (m *CafeRegistration) GetSig() []byte {
	if m != nil {
		return m.Sig
	}
	return nil
}

type CafeSession struct {
	Access               string               `protobuf:"bytes,1,opt,name=access,proto3" json:"access,omitempty"`
	Exp                  *timestamp.Timestamp `protobuf:"bytes,2,opt,name=exp,proto3" json:"exp,omitempty"`
	Refresh              string               `protobuf:"bytes,3,opt,name=refresh,proto3" json:"refresh,omitempty"`
	Rexp                 *timestamp.Timestamp `protobuf:"bytes,4,opt,name=rexp,proto3" json:"rexp,omitempty"`
	Subject              string               `protobuf:"bytes,5,opt,name=subject,proto3" json:"subject,omitempty"`
	Type                 string               `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CafeSession) Reset()         { *m = CafeSession{} }
func (m *CafeSession) String() string { return proto.CompactTextString(m) }
func (*CafeSession) ProtoMessage()    {}
func (*CafeSession) Descriptor() ([]byte, []int) {
	return fileDescriptor_cafe_93e3fd513c792df6, []int{3}
}
func (m *CafeSession) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CafeSession.Unmarshal(m, b)
}
func (m *CafeSession) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CafeSession.Marshal(b, m, deterministic)
}
func (dst *CafeSession) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CafeSession.Merge(dst, src)
}
func (m *CafeSession) XXX_Size() int {
	return xxx_messageInfo_CafeSession.Size(m)
}
func (m *CafeSession) XXX_DiscardUnknown() {
	xxx_messageInfo_CafeSession.DiscardUnknown(m)
}

var xxx_messageInfo_CafeSession proto.InternalMessageInfo

func (m *CafeSession) GetAccess() string {
	if m != nil {
		return m.Access
	}
	return ""
}

func (m *CafeSession) GetExp() *timestamp.Timestamp {
	if m != nil {
		return m.Exp
	}
	return nil
}

func (m *CafeSession) GetRefresh() string {
	if m != nil {
		return m.Refresh
	}
	return ""
}

func (m *CafeSession) GetRexp() *timestamp.Timestamp {
	if m != nil {
		return m.Rexp
	}
	return nil
}

func (m *CafeSession) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *CafeSession) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type CafeCidList struct {
	Cids                 []string `protobuf:"bytes,1,rep,name=cids,proto3" json:"cids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CafeCidList) Reset()         { *m = CafeCidList{} }
func (m *CafeCidList) String() string { return proto.CompactTextString(m) }
func (*CafeCidList) ProtoMessage()    {}
func (*CafeCidList) Descriptor() ([]byte, []int) {
	return fileDescriptor_cafe_93e3fd513c792df6, []int{4}
}
func (m *CafeCidList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CafeCidList.Unmarshal(m, b)
}
func (m *CafeCidList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CafeCidList.Marshal(b, m, deterministic)
}
func (dst *CafeCidList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CafeCidList.Merge(dst, src)
}
func (m *CafeCidList) XXX_Size() int {
	return xxx_messageInfo_CafeCidList.Size(m)
}
func (m *CafeCidList) XXX_DiscardUnknown() {
	xxx_messageInfo_CafeCidList.DiscardUnknown(m)
}

var xxx_messageInfo_CafeCidList proto.InternalMessageInfo

func (m *CafeCidList) GetCids() []string {
	if m != nil {
		return m.Cids
	}
	return nil
}

type CafeBlock struct {
	RawData              []byte   `protobuf:"bytes,1,opt,name=rawData,proto3" json:"rawData,omitempty"`
	Cid                  string   `protobuf:"bytes,2,opt,name=cid,proto3" json:"cid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CafeBlock) Reset()         { *m = CafeBlock{} }
func (m *CafeBlock) String() string { return proto.CompactTextString(m) }
func (*CafeBlock) ProtoMessage()    {}
func (*CafeBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_cafe_93e3fd513c792df6, []int{5}
}
func (m *CafeBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CafeBlock.Unmarshal(m, b)
}
func (m *CafeBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CafeBlock.Marshal(b, m, deterministic)
}
func (dst *CafeBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CafeBlock.Merge(dst, src)
}
func (m *CafeBlock) XXX_Size() int {
	return xxx_messageInfo_CafeBlock.Size(m)
}
func (m *CafeBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_CafeBlock.DiscardUnknown(m)
}

var xxx_messageInfo_CafeBlock proto.InternalMessageInfo

func (m *CafeBlock) GetRawData() []byte {
	if m != nil {
		return m.RawData
	}
	return nil
}

func (m *CafeBlock) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

func init() {
	proto.RegisterType((*CafeChallenge)(nil), "CafeChallenge")
	proto.RegisterType((*CafeNonce)(nil), "CafeNonce")
	proto.RegisterType((*CafeRegistration)(nil), "CafeRegistration")
	proto.RegisterType((*CafeSession)(nil), "CafeSession")
	proto.RegisterType((*CafeCidList)(nil), "CafeCidList")
	proto.RegisterType((*CafeBlock)(nil), "CafeBlock")
}

func init() { proto.RegisterFile("cafe.proto", fileDescriptor_cafe_93e3fd513c792df6) }

var fileDescriptor_cafe_93e3fd513c792df6 = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0x49, 0x93, 0x56, 0x3a, 0xad, 0x50, 0x16, 0x91, 0xd0, 0x8b, 0x6d, 0x4e, 0x15, 0x24,
	0x05, 0x3d, 0x78, 0x6f, 0x3d, 0x8a, 0x87, 0xe8, 0xc9, 0xdb, 0x66, 0x33, 0x49, 0xb7, 0xa6, 0xd9,
	0xb0, 0xbb, 0xf5, 0xcf, 0x67, 0xf4, 0x4b, 0xc9, 0xec, 0x26, 0xd4, 0x8b, 0x78, 0x9b, 0x37, 0xfc,
	0x78, 0x6f, 0xe6, 0x01, 0x08, 0x5e, 0x62, 0xda, 0x6a, 0x65, 0xd5, 0xfc, 0xaa, 0x52, 0xaa, 0xaa,
	0x71, 0xed, 0x54, 0x7e, 0x2c, 0xd7, 0x56, 0x1e, 0xd0, 0x58, 0x7e, 0x68, 0x3d, 0x90, 0x5c, 0xc3,
	0xf9, 0x96, 0x97, 0xb8, 0xdd, 0xf1, 0xba, 0xc6, 0xa6, 0x42, 0x16, 0xc3, 0x19, 0x2f, 0x0a, 0x8d,
	0xc6, 0xc4, 0xc1, 0x22, 0x58, 0x8d, 0xb3, 0x5e, 0x26, 0x4b, 0x18, 0x13, 0xfa, 0xa4, 0x1a, 0x81,
	0xec, 0x02, 0x86, 0xef, 0xbc, 0x3e, 0x62, 0x07, 0x79, 0x91, 0xec, 0x61, 0x46, 0x48, 0x86, 0x95,
	0x34, 0x56, 0x73, 0x2b, 0x55, 0xf3, 0xb7, 0xe1, 0xc9, 0x63, 0xf0, 0xcb, 0x83, 0xb6, 0x0d, 0x45,
	0xc4, 0xa1, 0xdf, 0x3a, 0xc1, 0x66, 0x10, 0x1a, 0x59, 0xc5, 0xd1, 0x22, 0x58, 0x4d, 0x33, 0x1a,
	0x93, 0xef, 0x00, 0x26, 0x14, 0xf6, 0x8c, 0xc6, 0x50, 0xce, 0x25, 0x8c, 0xb8, 0x10, 0xa7, 0x98,
	0x4e, 0xb1, 0x1b, 0x08, 0xf1, 0xb3, 0x75, 0x19, 0x93, 0xdb, 0x79, 0xea, 0x0b, 0x49, 0xfb, 0x42,
	0xd2, 0x97, 0xbe, 0x90, 0x8c, 0x30, 0xba, 0x56, 0x63, 0xa9, 0xd1, 0xec, 0xba, 0xfc, 0x5e, 0xb2,
	0x14, 0x22, 0x4d, 0x46, 0xd1, 0xbf, 0x46, 0x8e, 0x23, 0x27, 0x73, 0xcc, 0xf7, 0x28, 0x6c, 0x3c,
	0xf4, 0x4e, 0x9d, 0x64, 0x0c, 0x22, 0xfb, 0xd5, 0x62, 0x3c, 0x72, 0x6b, 0x37, 0x27, 0x4b, 0xff,
	0xcc, 0x56, 0x16, 0x8f, 0xd2, 0x38, 0x44, 0xc8, 0x82, 0x5e, 0x09, 0x09, 0xa1, 0x39, 0xb9, 0xf7,
	0xfd, 0x6f, 0x6a, 0x25, 0xde, 0xdc, 0x9d, 0xfc, 0xe3, 0x81, 0x5b, 0xee, 0xde, 0x9d, 0x66, 0xbd,
	0xa4, 0xa6, 0x84, 0x2c, 0xba, 0x4e, 0x69, 0xdc, 0x44, 0xaf, 0x83, 0x36, 0xcf, 0x47, 0xee, 0xd2,
	0xbb, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x45, 0xf8, 0x60, 0xdb, 0x1f, 0x02, 0x00, 0x00,
}
