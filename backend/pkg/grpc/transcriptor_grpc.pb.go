// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: transcriptor.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TranscriptorServiceClient is the client API for TranscriptorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TranscriptorServiceClient interface {
	StreamWav(ctx context.Context, opts ...grpc.CallOption) (TranscriptorService_StreamWavClient, error)
}

type transcriptorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTranscriptorServiceClient(cc grpc.ClientConnInterface) TranscriptorServiceClient {
	return &transcriptorServiceClient{cc}
}

func (c *transcriptorServiceClient) StreamWav(ctx context.Context, opts ...grpc.CallOption) (TranscriptorService_StreamWavClient, error) {
	stream, err := c.cc.NewStream(ctx, &TranscriptorService_ServiceDesc.Streams[0], "/goMeetingTranscriptor.TranscriptorService/StreamWav", opts...)
	if err != nil {
		return nil, err
	}
	x := &transcriptorServiceStreamWavClient{stream}
	return x, nil
}

type TranscriptorService_StreamWavClient interface {
	Send(*WavChunk) error
	CloseAndRecv() (*WavResponse, error)
	grpc.ClientStream
}

type transcriptorServiceStreamWavClient struct {
	grpc.ClientStream
}

func (x *transcriptorServiceStreamWavClient) Send(m *WavChunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *transcriptorServiceStreamWavClient) CloseAndRecv() (*WavResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(WavResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TranscriptorServiceServer is the server API for TranscriptorService service.
// All implementations must embed UnimplementedTranscriptorServiceServer
// for forward compatibility
type TranscriptorServiceServer interface {
	StreamWav(TranscriptorService_StreamWavServer) error
	mustEmbedUnimplementedTranscriptorServiceServer()
}

// UnimplementedTranscriptorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTranscriptorServiceServer struct {
}

func (UnimplementedTranscriptorServiceServer) StreamWav(TranscriptorService_StreamWavServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamWav not implemented")
}
func (UnimplementedTranscriptorServiceServer) mustEmbedUnimplementedTranscriptorServiceServer() {}

// UnsafeTranscriptorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TranscriptorServiceServer will
// result in compilation errors.
type UnsafeTranscriptorServiceServer interface {
	mustEmbedUnimplementedTranscriptorServiceServer()
}

func RegisterTranscriptorServiceServer(s grpc.ServiceRegistrar, srv TranscriptorServiceServer) {
	s.RegisterService(&TranscriptorService_ServiceDesc, srv)
}

func _TranscriptorService_StreamWav_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TranscriptorServiceServer).StreamWav(&transcriptorServiceStreamWavServer{stream})
}

type TranscriptorService_StreamWavServer interface {
	SendAndClose(*WavResponse) error
	Recv() (*WavChunk, error)
	grpc.ServerStream
}

type transcriptorServiceStreamWavServer struct {
	grpc.ServerStream
}

func (x *transcriptorServiceStreamWavServer) SendAndClose(m *WavResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *transcriptorServiceStreamWavServer) Recv() (*WavChunk, error) {
	m := new(WavChunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TranscriptorService_ServiceDesc is the grpc.ServiceDesc for TranscriptorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TranscriptorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "goMeetingTranscriptor.TranscriptorService",
	HandlerType: (*TranscriptorServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamWav",
			Handler:       _TranscriptorService_StreamWav_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "transcriptor.proto",
}
