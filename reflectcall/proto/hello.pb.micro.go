package proto

import (
	"context"
	"github.com/micro/go-micro/v2/server"
)

type HelloHandler interface {
	Call(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, Hello_StreamStream) error
	PingPong(context.Context, Hello_PingPongStream) error
}

type Hello_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type Hello_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

func RegisterHelloHandler(s server.Server, hdlr HelloHandler, opts ...server.HandlerOption) error {
	type hello interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Hello struct {
		hello
	}
	h := &helloHandler{hdlr}

	return s.Handle(s.NewHandler(&Hello{h}, opts...))
}

type helloHandler struct {
	HelloHandler
}

func (h *helloHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.HelloHandler.Call(ctx, in, out)
}

func (h *helloHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.HelloHandler.Stream(ctx, m, &helloStreamStream{stream})
}

func (h *helloHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.HelloHandler.PingPong(ctx, &helloPingPongStream{stream})
}

type helloStreamStream struct {
	stream server.Stream
}

func (x *helloStreamStream) Close() error {
	return x.stream.Close()
}

func (x *helloStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *helloStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *helloStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *helloStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

type helloPingPongStream struct {
	stream server.Stream
}

func (x *helloPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *helloPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *helloPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *helloPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *helloPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *helloPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
