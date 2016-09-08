package multiinterceptor

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type multiUnaryServerInterceptor struct {
	uints []grpc.UnaryServerInterceptor
}

func NewMultiUnaryServerInterceptor(uints ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return (&multiUnaryServerInterceptor{uints: uints}).gen
}

func (m *multiUnaryServerInterceptor) gen(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return m.chain(0, ctx, req, info, handler)
}

func (m *multiUnaryServerInterceptor) chain(i int, ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if i == len(m.uints) {
		return handler(ctx, req)
	}
	return m.uints[i](ctx, req, info, func(ctx2 context.Context, req2 interface{}) (interface{}, error) {
		return m.chain(i+1, ctx2, req2, info, handler)
	})
}

type multiStreamServerInterceptor struct {
	sints []grpc.StreamServerInterceptor
}

func NewMultiStreamServerInterceptor(sints ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return (&multiStreamServerInterceptor{sints: sints}).gen
}

func (m *multiStreamServerInterceptor) gen(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return m.chain(0, srv, stream, info, handler)
}

func (m *multiStreamServerInterceptor) chain(i int, srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if i == len(m.sints) {
		return handler(srv, stream)
	}
	return m.sints[i](srv, stream, info, func(srv2 interface{}, stream2 grpc.ServerStream) error {
		return m.chain(i+1, srv2, stream2, info, handler)
	})
}

type serverStreamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss serverStreamWithContext) Context() context.Context {
	return ss.ctx
}

func NewServerStreamWithContext(stream grpc.ServerStream, ctx context.Context) grpc.ServerStream {
	return serverStreamWithContext{
		ServerStream: stream,
		ctx:          ctx,
	}
}
