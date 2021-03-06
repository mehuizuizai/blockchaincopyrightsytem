// Code generated by protoc-gen-go. DO NOT EDIT.
// source: consensus.proto

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// consensus request body
type ConsensusRequest struct {
	SessionID string `protobuf:"bytes,1,opt,name=SessionID" json:"SessionID,omitempty"`
	IP        string `protobuf:"bytes,2,opt,name=IP" json:"IP,omitempty"`
	Vote      []byte `protobuf:"bytes,3,opt,name=Vote,proto3" json:"Vote,omitempty"`
}

func (m *ConsensusRequest) Reset()                    { *m = ConsensusRequest{} }
func (m *ConsensusRequest) String() string            { return proto1.CompactTextString(m) }
func (*ConsensusRequest) ProtoMessage()               {}
func (*ConsensusRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *ConsensusRequest) GetSessionID() string {
	if m != nil {
		return m.SessionID
	}
	return ""
}

func (m *ConsensusRequest) GetIP() string {
	if m != nil {
		return m.IP
	}
	return ""
}

func (m *ConsensusRequest) GetVote() []byte {
	if m != nil {
		return m.Vote
	}
	return nil
}

// consensus response body
type ConsensusResponse struct {
}

func (m *ConsensusResponse) Reset()                    { *m = ConsensusResponse{} }
func (m *ConsensusResponse) String() string            { return proto1.CompactTextString(m) }
func (*ConsensusResponse) ProtoMessage()               {}
func (*ConsensusResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func init() {
	proto1.RegisterType((*ConsensusRequest)(nil), "proto.ConsensusRequest")
	proto1.RegisterType((*ConsensusResponse)(nil), "proto.ConsensusResponse")
}

func init() { proto1.RegisterFile("consensus.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 125 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0xce, 0xcf, 0x2b,
	0x4e, 0xcd, 0x2b, 0x2e, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a,
	0x21, 0x5c, 0x02, 0xce, 0x30, 0x99, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x19, 0x2e,
	0xce, 0xe0, 0xd4, 0xe2, 0xe2, 0xcc, 0xfc, 0x3c, 0x4f, 0x17, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce,
	0x20, 0x84, 0x80, 0x10, 0x1f, 0x17, 0x93, 0x67, 0x80, 0x04, 0x13, 0x58, 0x98, 0xc9, 0x33, 0x40,
	0x48, 0x88, 0x8b, 0x25, 0x2c, 0xbf, 0x24, 0x55, 0x82, 0x59, 0x81, 0x51, 0x83, 0x27, 0x08, 0xcc,
	0x56, 0x12, 0xe6, 0x12, 0x44, 0x32, 0xb5, 0xb8, 0x00, 0xc4, 0x4e, 0x62, 0x03, 0xdb, 0x68, 0x0c,
	0x08, 0x00, 0x00, 0xff, 0xff, 0xa3, 0x82, 0x64, 0xa7, 0x8b, 0x00, 0x00, 0x00,
}
