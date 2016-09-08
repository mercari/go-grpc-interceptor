package instrument

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var ()

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
	callCounter := make(map[string]int)
	callErrors := 0
	InstallInstrumentor(func(call MethodCall) {
		callCounter[call.FullMethod]++
		if call.Error != nil {
			callErrors++
		}
	})

	var count = 0
	methodName := "TestService.UnaryMethod"
	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: methodName,
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		if count == 1 {
			return nil, fmt.Errorf("error")
		}
		return "output", nil
	}

	for ; count < 3; count++ {
		ctx := context.Background()
		UnaryServerInterceptor(ctx, "xyz", unaryInfo, unaryHandler)
	}

	if got, want := callCounter[methodName], 3; got != want {
		t.Errorf("expect counter to %q, but got %q", want, got)
	}
	if got, want := callErrors, 1; got != want {
		t.Errorf("expect errors to %q, but got %q", want, got)
	}
}

func TestStreamServer(t *testing.T) {
	callCounter := make(map[string]int)
	callErrors := 0
	InstallInstrumentor(func(call MethodCall) {
		callCounter[call.FullMethod]++
		if call.Error != nil {
			callErrors++
		}
	})

	var count = 0
	methodName := "TestService.StreamMethod"
	streamInfo := &grpc.StreamServerInfo{
		FullMethod:     methodName,
		IsServerStream: true,
	}
	streamHandler := func(srv interface{}, stream grpc.ServerStream) error {
		if count == 1 {
			return fmt.Errorf("error")
		}
		return nil
	}
	testService := struct{}{}
	testStream := &testServerStream{ctx: context.Background()}

	for ; count < 3; count++ {
		StreamServerInterceptor(testService, testStream, streamInfo, streamHandler)
	}

	if got, want := callCounter[methodName], 3; got != want {
		t.Errorf("expect counter to %q, but got %q", want, got)
	}
	if got, want := callErrors, 1; got != want {
		t.Errorf("expect errors to %q, but got %q", want, got)
	}
}
