package panichandler

import (
	"strings"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	doPanic := true

	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		if doPanic {
			panic("test error")
		}
		return "output", nil
	}

	ctx := context.Background()
	_, err := UnaryServerInterceptor(ctx, "xyz", unaryInfo, unaryHandler)
	if err == nil {
		t.Fatalf("unexpected success")
	}

	if got, want := grpc.Code(err), codes.Internal; got != want {
		t.Errorf("expect grpc.Code to %s, but got %s", want, got)
	}
	if got := grpc.ErrorDesc(err); !strings.HasPrefix(got, "panic") {
		t.Errorf("expect ErrorDesc has %q prefix: %q", "panic", got)
	}

	doPanic = false
	out, err := UnaryServerInterceptor(ctx, "xyz", unaryInfo, unaryHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := out, "output"; got != want {
		t.Errorf("expect output to %s, but got %s", want, got)
	}
}

func TestStreamServer(t *testing.T) {
	doPanic := true

	streamInfo := &grpc.StreamServerInfo{
		FullMethod:     "TestService.StreamMethod",
		IsServerStream: true,
	}
	streamHandler := func(srv interface{}, stream grpc.ServerStream) error {
		if doPanic {
			panic("test error")
		}
		return nil
	}
	testService := struct{}{}
	testStream := &testServerStream{ctx: context.Background()}

	err := StreamServerInterceptor(testService, testStream, streamInfo, streamHandler)
	if err == nil {
		t.Fatalf("unexpected success")
	}

	doPanic = false
	err = StreamServerInterceptor(testService, testStream, streamInfo, streamHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPanicDump(t *testing.T) {
	doPanic := true

	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		if doPanic {
			panic("test error")
		}
		return "output", nil
	}

	InstallPanicHandler(LogPanicDump)

	ctx := context.Background()
	_, err := UnaryServerInterceptor(ctx, "xyz", unaryInfo, unaryHandler)
	if err == nil {
		t.Fatalf("unexpected success")
	}

	if got, want := grpc.Code(err), codes.Internal; got != want {
		t.Errorf("expect grpc.Code to %s, but got %s", want, got)
	}
	if got := grpc.ErrorDesc(err); !strings.HasPrefix(got, "panic") {
		t.Errorf("expect ErrorDesc has %q prefix: %q", "panic", got)
	}
}
