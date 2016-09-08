package zap

import (
	"bytes"
	"testing"

	"github.com/mercari/go-grpc-interceptor/zap/zapctx"
	"github.com/uber-go/zap"
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
	buf := &bytes.Buffer{}
	encoder := zap.NewJSONEncoder(zap.NoTime())
	logger := zap.New(encoder, zap.DebugLevel, zap.Output(zap.AddSync(buf)), zap.ErrorOutput(zap.AddSync(buf)))

	expected := `{"level":"info","msg":"message","method":"TestService.UnaryMethod"}` + "\n"

	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		zapctx.MustFromContext(ctx).Info("message")

		if got, want := string(buf.Bytes()), expected; got != want {
			t.Errorf("\nexpected:\n%sgot:\n%s", want, got)
		}
		return "output", nil
	}

	ctx := context.Background()
	_, err := UnaryServerInterceptor(logger)(ctx, "xyz", unaryInfo, unaryHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStreamServer(t *testing.T) {
	buf := &bytes.Buffer{}
	encoder := zap.NewJSONEncoder(zap.NoTime())
	logger := zap.New(encoder, zap.DebugLevel, zap.Output(zap.AddSync(buf)), zap.ErrorOutput(zap.AddSync(buf)))

	expected := `{"level":"info","msg":"message","method":"TestService.StreamMethod"}` + "\n"

	streamInfo := &grpc.StreamServerInfo{
		FullMethod:     "TestService.StreamMethod",
		IsServerStream: true,
	}
	streamHandler := func(srv interface{}, stream grpc.ServerStream) error {
		zapctx.MustFromContext(stream.Context()).Info("message")

		if got, want := string(buf.Bytes()), expected; got != want {
			t.Errorf("\nexpected:\n%sgot:\n%s", want, got)
		}

		return nil
	}
	testService := struct{}{}
	testStream := &testServerStream{ctx: context.Background()}

	err := StreamServerInterceptor(logger)(testService, testStream, streamInfo, streamHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
