// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chat.proto

package grpc

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ChatStream_StreamType int32

const (
	ChatStream_JOINED ChatStream_StreamType = 0
	ChatStream_LEAVE  ChatStream_StreamType = 1
	ChatStream_LEAVED ChatStream_StreamType = 2
	ChatStream_SPEAK  ChatStream_StreamType = 3
	ChatStream_SPOKE  ChatStream_StreamType = 4
)

var ChatStream_StreamType_name = map[int32]string{
	0: "JOINED",
	1: "LEAVE",
	2: "LEAVED",
	3: "SPEAK",
	4: "SPOKE",
}

var ChatStream_StreamType_value = map[string]int32{
	"JOINED": 0,
	"LEAVE":  1,
	"LEAVED": 2,
	"SPEAK":  3,
	"SPOKE":  4,
}

func (x ChatStream_StreamType) String() string {
	return proto.EnumName(ChatStream_StreamType_name, int32(x))
}

func (ChatStream_StreamType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{4, 0}
}

type JoinRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JoinRequest) Reset()         { *m = JoinRequest{} }
func (m *JoinRequest) String() string { return proto.CompactTextString(m) }
func (*JoinRequest) ProtoMessage()    {}
func (*JoinRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{0}
}

func (m *JoinRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinRequest.Unmarshal(m, b)
}
func (m *JoinRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinRequest.Marshal(b, m, deterministic)
}
func (m *JoinRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinRequest.Merge(m, src)
}
func (m *JoinRequest) XXX_Size() int {
	return xxx_messageInfo_JoinRequest.Size(m)
}
func (m *JoinRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinRequest.DiscardUnknown(m)
}

var xxx_messageInfo_JoinRequest proto.InternalMessageInfo

func (m *JoinRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *JoinRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type LeaveRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LeaveRequest) Reset()         { *m = LeaveRequest{} }
func (m *LeaveRequest) String() string { return proto.CompactTextString(m) }
func (*LeaveRequest) ProtoMessage()    {}
func (*LeaveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{1}
}

func (m *LeaveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeaveRequest.Unmarshal(m, b)
}
func (m *LeaveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeaveRequest.Marshal(b, m, deterministic)
}
func (m *LeaveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaveRequest.Merge(m, src)
}
func (m *LeaveRequest) XXX_Size() int {
	return xxx_messageInfo_LeaveRequest.Size(m)
}
func (m *LeaveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LeaveRequest proto.InternalMessageInfo

func (m *LeaveRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type SpeakRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SpeakRequest) Reset()         { *m = SpeakRequest{} }
func (m *SpeakRequest) String() string { return proto.CompactTextString(m) }
func (*SpeakRequest) ProtoMessage()    {}
func (*SpeakRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{2}
}

func (m *SpeakRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpeakRequest.Unmarshal(m, b)
}
func (m *SpeakRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpeakRequest.Marshal(b, m, deterministic)
}
func (m *SpeakRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpeakRequest.Merge(m, src)
}
func (m *SpeakRequest) XXX_Size() int {
	return xxx_messageInfo_SpeakRequest.Size(m)
}
func (m *SpeakRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SpeakRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SpeakRequest proto.InternalMessageInfo

func (m *SpeakRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *SpeakRequest) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type CommonResponse struct {
	Result               bool     `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommonResponse) Reset()         { *m = CommonResponse{} }
func (m *CommonResponse) String() string { return proto.CompactTextString(m) }
func (*CommonResponse) ProtoMessage()    {}
func (*CommonResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{3}
}

func (m *CommonResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonResponse.Unmarshal(m, b)
}
func (m *CommonResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonResponse.Marshal(b, m, deterministic)
}
func (m *CommonResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonResponse.Merge(m, src)
}
func (m *CommonResponse) XXX_Size() int {
	return xxx_messageInfo_CommonResponse.Size(m)
}
func (m *CommonResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CommonResponse proto.InternalMessageInfo

func (m *CommonResponse) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

type ChatStream struct {
	Type                 ChatStream_StreamType `protobuf:"varint,1,opt,name=type,proto3,enum=grpc.ChatStream_StreamType" json:"type,omitempty"`
	Uuid                 string                `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Name                 string                `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Msg                  string                `protobuf:"bytes,4,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ChatStream) Reset()         { *m = ChatStream{} }
func (m *ChatStream) String() string { return proto.CompactTextString(m) }
func (*ChatStream) ProtoMessage()    {}
func (*ChatStream) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{4}
}

func (m *ChatStream) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChatStream.Unmarshal(m, b)
}
func (m *ChatStream) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChatStream.Marshal(b, m, deterministic)
}
func (m *ChatStream) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChatStream.Merge(m, src)
}
func (m *ChatStream) XXX_Size() int {
	return xxx_messageInfo_ChatStream.Size(m)
}
func (m *ChatStream) XXX_DiscardUnknown() {
	xxx_messageInfo_ChatStream.DiscardUnknown(m)
}

var xxx_messageInfo_ChatStream proto.InternalMessageInfo

func (m *ChatStream) GetType() ChatStream_StreamType {
	if m != nil {
		return m.Type
	}
	return ChatStream_JOINED
}

func (m *ChatStream) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *ChatStream) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ChatStream) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterEnum("grpc.ChatStream_StreamType", ChatStream_StreamType_name, ChatStream_StreamType_value)
	proto.RegisterType((*JoinRequest)(nil), "grpc.JoinRequest")
	proto.RegisterType((*LeaveRequest)(nil), "grpc.LeaveRequest")
	proto.RegisterType((*SpeakRequest)(nil), "grpc.SpeakRequest")
	proto.RegisterType((*CommonResponse)(nil), "grpc.CommonResponse")
	proto.RegisterType((*ChatStream)(nil), "grpc.ChatStream")
}

func init() { proto.RegisterFile("chat.proto", fileDescriptor_8c585a45e2093e54) }

var fileDescriptor_8c585a45e2093e54 = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0xd1, 0x4e, 0xc2, 0x30,
	0x14, 0xa5, 0x50, 0x88, 0x5c, 0x09, 0xa9, 0x8d, 0x31, 0x04, 0x5f, 0x4c, 0x9f, 0x78, 0x9a, 0x06,
	0xf4, 0x03, 0x08, 0xf4, 0x41, 0x20, 0x62, 0x86, 0xf1, 0xbd, 0xe2, 0x0d, 0x10, 0xdd, 0x5a, 0xb7,
	0xce, 0x84, 0xcf, 0xf1, 0x47, 0xfc, 0x36, 0xd3, 0x0e, 0xdc, 0x34, 0xca, 0xd3, 0xce, 0x4e, 0x7b,
	0x6e, 0xce, 0xb9, 0xa7, 0x00, 0xcb, 0xb5, 0xb2, 0x81, 0x49, 0xb4, 0xd5, 0x9c, 0xae, 0x12, 0xb3,
	0x14, 0x37, 0x70, 0x3c, 0xd1, 0x9b, 0x38, 0xc4, 0xb7, 0x0c, 0x53, 0xcb, 0x39, 0xd0, 0x2c, 0xdb,
	0x3c, 0x77, 0xc8, 0x05, 0xe9, 0x35, 0x43, 0x8f, 0x1d, 0x17, 0xab, 0x08, 0x3b, 0xd5, 0x9c, 0x73,
	0x58, 0x08, 0x68, 0xcd, 0x50, 0xbd, 0xe3, 0x01, 0x9d, 0xb8, 0x86, 0xd6, 0xc2, 0xa0, 0x7a, 0x39,
	0x34, 0x9b, 0x41, 0x2d, 0x4a, 0x57, 0xbb, 0xd1, 0x0e, 0x8a, 0x1e, 0xb4, 0x47, 0x3a, 0x8a, 0x74,
	0x1c, 0x62, 0x6a, 0x74, 0x9c, 0x22, 0x3f, 0x83, 0x46, 0x82, 0x69, 0xf6, 0x6a, 0xbd, 0xf2, 0x28,
	0xdc, 0xfd, 0x89, 0x4f, 0x02, 0x30, 0x5a, 0x2b, 0xbb, 0xb0, 0x09, 0xaa, 0x88, 0x5f, 0x02, 0xb5,
	0x5b, 0x83, 0xfe, 0x52, 0xbb, 0x7f, 0x1e, 0xb8, 0x78, 0x41, 0x71, 0x1e, 0xe4, 0x9f, 0x87, 0xad,
	0xc1, 0xd0, 0x5f, 0xfc, 0xf6, 0x53, 0xfd, 0x23, 0x6b, 0xad, 0xc8, 0xba, 0xf7, 0x48, 0x0b, 0x8f,
	0x12, 0xa0, 0x98, 0xc6, 0x01, 0x1a, 0x93, 0xf9, 0xed, 0x9d, 0x1c, 0xb3, 0x0a, 0x6f, 0x42, 0x7d,
	0x26, 0x87, 0x8f, 0x92, 0x11, 0x47, 0x7b, 0x38, 0x66, 0x55, 0x47, 0x2f, 0xee, 0xe5, 0x70, 0xca,
	0x6a, 0x39, 0x9c, 0x4f, 0x25, 0xa3, 0xfd, 0x0f, 0x02, 0xd4, 0x19, 0x74, 0xd6, 0x5d, 0x09, 0xfc,
	0x24, 0x37, 0x5d, 0x2a, 0xa4, 0xcb, 0x7e, 0xe7, 0x10, 0x95, 0x2b, 0xc2, 0x07, 0x50, 0xf7, 0xeb,
	0xe7, 0x3c, 0x3f, 0x2e, 0x77, 0xd1, 0x3d, 0xdd, 0x49, 0x7e, 0x6c, 0x51, 0x54, 0x9c, 0xc8, 0xf7,
	0xb1, 0x17, 0x95, 0xcb, 0xf9, 0x4f, 0xf4, 0xd4, 0xf0, 0x8f, 0x65, 0xf0, 0x15, 0x00, 0x00, 0xff,
	0xff, 0x63, 0x61, 0xc6, 0xb3, 0x3a, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChatClient interface {
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (Chat_JoinClient, error)
	Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*CommonResponse, error)
	Speak(ctx context.Context, in *SpeakRequest, opts ...grpc.CallOption) (*CommonResponse, error)
}

type chatClient struct {
	cc *grpc.ClientConn
}

func NewChatClient(cc *grpc.ClientConn) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (Chat_JoinClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Chat_serviceDesc.Streams[0], "/grpc.Chat/Join", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatJoinClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Chat_JoinClient interface {
	Recv() (*ChatStream, error)
	grpc.ClientStream
}

type chatJoinClient struct {
	grpc.ClientStream
}

func (x *chatJoinClient) Recv() (*ChatStream, error) {
	m := new(ChatStream)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatClient) Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, "/grpc.Chat/Leave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) Speak(ctx context.Context, in *SpeakRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, "/grpc.Chat/Speak", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
type ChatServer interface {
	Join(*JoinRequest, Chat_JoinServer) error
	Leave(context.Context, *LeaveRequest) (*CommonResponse, error)
	Speak(context.Context, *SpeakRequest) (*CommonResponse, error)
}

// UnimplementedChatServer can be embedded to have forward compatible implementations.
type UnimplementedChatServer struct {
}

func (*UnimplementedChatServer) Join(req *JoinRequest, srv Chat_JoinServer) error {
	return status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (*UnimplementedChatServer) Leave(ctx context.Context, req *LeaveRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (*UnimplementedChatServer) Speak(ctx context.Context, req *SpeakRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Speak not implemented")
}

func RegisterChatServer(s *grpc.Server, srv ChatServer) {
	s.RegisterService(&_Chat_serviceDesc, srv)
}

func _Chat_Join_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(JoinRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServer).Join(m, &chatJoinServer{stream})
}

type Chat_JoinServer interface {
	Send(*ChatStream) error
	grpc.ServerStream
}

type chatJoinServer struct {
	grpc.ServerStream
}

func (x *chatJoinServer) Send(m *ChatStream) error {
	return x.ServerStream.SendMsg(m)
}

func _Chat_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Chat/Leave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).Leave(ctx, req.(*LeaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_Speak_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SpeakRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).Speak(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Chat/Speak",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).Speak(ctx, req.(*SpeakRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Chat_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Leave",
			Handler:    _Chat_Leave_Handler,
		},
		{
			MethodName: "Speak",
			Handler:    _Chat_Speak_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Join",
			Handler:       _Chat_Join_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chat.proto",
}