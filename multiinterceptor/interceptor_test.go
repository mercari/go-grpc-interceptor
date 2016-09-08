package multiinterceptor

import (
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type testServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *testServerStream) Context() context.Context {
	return ss.ctx
}

func (ss *testServerStream) SendMsg(m interface{}) error {
	return nil
}

func (f *testServerStream) RecvMsg(m interface{}) error {
	return nil
}

func TestUnaryServer(t *testing.T) {
	input := "input"

	foo := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if got, want := req.(string), input; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}

		ctx = context.WithValue(ctx, "foo", 1)
		return handler(ctx, req)
	}
	bar := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if got, want := req.(string), input; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}
		if got, want := ctx.Value("foo").(int), 1; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}

		ctx = context.WithValue(ctx, "bar", 2)
		return handler(ctx, req)
	}

	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		if got, want := req.(string), input; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}
		if got, want := ctx.Value("foo").(int), 1; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}
		if got, want := ctx.Value("bar").(int), 2; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}

		return "output", nil
	}

	ctx := context.Background()
	interceptor := NewMultiUnaryServerInterceptor(foo, bar)
	out, err := interceptor(ctx, input, unaryInfo, unaryHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := out, "output"; got != want {
		t.Errorf("expect output to %q, but got %q", want, got)
	}
}

func TestStreamServer(t *testing.T) {
	foo := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		ctx = context.WithValue(ctx, "foo", 1)
		stream = NewServerStreamWithContext(stream, ctx)
		return handler(srv, stream)
	}
	bar := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		if got, want := ctx.Value("foo").(int), 1; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}

		ctx = context.WithValue(ctx, "bar", 2)
		stream = NewServerStreamWithContext(stream, ctx)
		return handler(srv, stream)
	}

	streamInfo := &grpc.StreamServerInfo{
		FullMethod:     "TestService.StreamMethod",
		IsServerStream: true,
	}
	streamHandler := func(srv interface{}, stream grpc.ServerStream) error {
		ctx := stream.Context()
		if got, want := ctx.Value("foo").(int), 1; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}
		if got, want := ctx.Value("bar").(int), 2; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}
		return nil
	}
	testService := struct{}{}
	testStream := &testServerStream{ctx: context.Background()}

	interceptor := NewMultiStreamServerInterceptor(foo, bar)
	err := interceptor(testService, testStream, streamInfo, streamHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
