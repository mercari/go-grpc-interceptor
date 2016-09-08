package instrument

import (
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type MethodCall struct {
	Context     context.Context
	StartedAt   time.Time
	Duration    time.Duration
	FullMethod  string
	Service     string
	Method      string
	ServerInfo  interface{}
	IsStreaming bool
	Error       error
}

type Instrumentor func(MethodCall)

var instrument Instrumentor = nullInstrument

func InstallInstrumentor(i Instrumentor) {
	instrument = i
}

var _ grpc.UnaryServerInterceptor = UnaryServerInterceptor
var _ grpc.StreamServerInterceptor = StreamServerInterceptor

func splitFullMethod(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	pos := strings.Index(fullMethod, "/")
	if pos < 0 {
		return "unknown", "unknown"
	}
	return fullMethod[:pos], fullMethod[pos+1:]
}

func newMethodCall(ctx context.Context, fullMethod string, info interface{}, streaming bool, start time.Time, err error) MethodCall {
	service, method := splitFullMethod(fullMethod)
	return MethodCall{
		Context:     ctx,
		StartedAt:   start,
		Duration:    time.Since(start),
		FullMethod:  fullMethod,
		Service:     service,
		Method:      method,
		ServerInfo:  info,
		IsStreaming: streaming,
		Error:       err,
	}
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func(start time.Time) {
		instrument(newMethodCall(ctx, info.FullMethod, info, false, start, err))
	}(time.Now())
	return handler(ctx, req)
}

func StreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func(start time.Time) {
		instrument(newMethodCall(stream.Context(), info.FullMethod, info, true, start, err))
	}(time.Now())
	return handler(srv, stream)
}

func nullInstrument(call MethodCall) {
}
