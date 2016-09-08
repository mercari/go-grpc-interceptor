package zap

import (
	"bytes"
	"testing"

	"github.com/mercari/go-grpc-interceptor/panichandler"
	"github.com/mercari/go-grpc-interceptor/zap/zapctx"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestUnaryServer(t *testing.T) {
	buf := &bytes.Buffer{}
	encoder := zap.NewJSONEncoder(zap.NoTime())
	logger := zap.New(encoder, zap.DebugLevel, zap.Output(zap.AddSync(buf)), zap.ErrorOutput(zap.AddSync(buf)))

	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test error")
	}

	panichandler.InstallPanicHandler(LogPanicWithStackTrace)

	ctx := context.Background()
	ctx = zapctx.NewContext(ctx, logger)
	_, err := panichandler.UnaryServerInterceptor(ctx, "xyz", unaryInfo, unaryHandler)
	if err == nil {
		t.Fatalf("unexpected success")
	}

	if got, want := grpc.Code(err), codes.Internal; got != want {
		t.Errorf("expect grpc.Code to %s, but got %s", want, got)
	}
}
