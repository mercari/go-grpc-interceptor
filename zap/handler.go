package zap

import (
	multiint "github.com/eltorocorp/go-grpc-request-id-interceptor/multiinterceptor"
	"github.com/eltorocorp/go-grpc-request-id-interceptor/xrequestid"
	"github.com/eltorocorp/go-grpc-request-id-interceptor/zap/zapctx"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var DefaultMethodKey = "method"
var DefaultRequestIDKey = "requestid"

func UnaryServerInterceptor(logger zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		l := logger.With(zap.String(DefaultMethodKey, info.FullMethod))
		ctx = zapctx.NewContext(ctx, l)
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(logger zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		l := logger.With(zap.String(DefaultMethodKey, info.FullMethod))
		ctx := zapctx.NewContext(stream.Context(), l)
		stream = multiint.NewServerStreamWithContext(stream, ctx)
		return handler(srv, stream)
	}
}

func UnaryServerInterceptorWithRequestID(logger zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		l := logger.With(
			zap.String(DefaultMethodKey, info.FullMethod),
			zap.String(DefaultRequestIDKey, xrequestid.FromContext(ctx)),
		)
		ctx = zapctx.NewContext(ctx, l)
		return handler(ctx, req)
	}
}

func StreamServerInterceptorWithRequestID(logger zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		l := logger.With(
			zap.String(DefaultMethodKey, info.FullMethod),
			zap.String(DefaultRequestIDKey, xrequestid.FromContext(stream.Context())),
		)
		ctx := zapctx.NewContext(stream.Context(), l)
		stream = multiint.NewServerStreamWithContext(stream, ctx)
		return handler(srv, stream)
	}
}
