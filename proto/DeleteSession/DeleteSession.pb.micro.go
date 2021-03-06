// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/DeleteSession/DeleteSession.proto

package go_micro_srv_DeleteSession

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for DeleteSession service

type DeleteSessionService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (DeleteSession_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (DeleteSession_PingPongService, error)
}

type deleteSessionService struct {
	c    client.Client
	name string
}

func NewDeleteSessionService(name string, c client.Client) DeleteSessionService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.srv.DeleteSession"
	}
	return &deleteSessionService{
		c:    c,
		name: name,
	}
}

func (c *deleteSessionService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "DeleteSession.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deleteSessionService) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (DeleteSession_StreamService, error) {
	req := c.c.NewRequest(c.name, "DeleteSession.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &deleteSessionServiceStream{stream}, nil
}

type DeleteSession_StreamService interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type deleteSessionServiceStream struct {
	stream client.Stream
}

func (x *deleteSessionServiceStream) Close() error {
	return x.stream.Close()
}

func (x *deleteSessionServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *deleteSessionServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *deleteSessionServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *deleteSessionService) PingPong(ctx context.Context, opts ...client.CallOption) (DeleteSession_PingPongService, error) {
	req := c.c.NewRequest(c.name, "DeleteSession.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &deleteSessionServicePingPong{stream}, nil
}

type DeleteSession_PingPongService interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type deleteSessionServicePingPong struct {
	stream client.Stream
}

func (x *deleteSessionServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *deleteSessionServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *deleteSessionServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *deleteSessionServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *deleteSessionServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for DeleteSession service

type DeleteSessionHandler interface {
	Call(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, DeleteSession_StreamStream) error
	PingPong(context.Context, DeleteSession_PingPongStream) error
}

func RegisterDeleteSessionHandler(s server.Server, hdlr DeleteSessionHandler, opts ...server.HandlerOption) error {
	type deleteSession interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type DeleteSession struct {
		deleteSession
	}
	h := &deleteSessionHandler{hdlr}
	return s.Handle(s.NewHandler(&DeleteSession{h}, opts...))
}

type deleteSessionHandler struct {
	DeleteSessionHandler
}

func (h *deleteSessionHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.DeleteSessionHandler.Call(ctx, in, out)
}

func (h *deleteSessionHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.DeleteSessionHandler.Stream(ctx, m, &deleteSessionStreamStream{stream})
}

type DeleteSession_StreamStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type deleteSessionStreamStream struct {
	stream server.Stream
}

func (x *deleteSessionStreamStream) Close() error {
	return x.stream.Close()
}

func (x *deleteSessionStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *deleteSessionStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *deleteSessionStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *deleteSessionHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.DeleteSessionHandler.PingPong(ctx, &deleteSessionPingPongStream{stream})
}

type DeleteSession_PingPongStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type deleteSessionPingPongStream struct {
	stream server.Stream
}

func (x *deleteSessionPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *deleteSessionPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *deleteSessionPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *deleteSessionPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *deleteSessionPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
