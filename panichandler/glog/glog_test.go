package glog

import (
	"flag"
	"testing"

	"github.com/mercari/go-grpc-interceptor/panichandler"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestPanicDump(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test error")
	}

	panichandler.InstallPanicHandler(LogPanicStackMultiLine)

	ctx := context.Background()
	_, err := panichandler.UnaryServerInterceptor(ctx, "xyz", unaryInfo, unaryHandler)
	if err == nil {
		t.Fatalf("unexpected success")
	}

	if got, want := grpc.Code(err), codes.Internal; got != want {
		t.Errorf("expect grpc.Code to %s, but got %s", want, got)
	}
}
