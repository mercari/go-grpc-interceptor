package panichandler

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var _ grpc.UnaryServerInterceptor = UnaryServerInterceptor
var _ grpc.StreamServerInterceptor = StreamServerInterceptor

func toPanicError(r interface{}) error {
	return grpc.Errorf(codes.Internal, "panic: %v", r)
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer handleCrash(ctx, func(ctx context.Context, r interface{}) {
		err = toPanicError(r)
	})

	return handler(ctx, req)
}

func StreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer handleCrash(stream.Context(), func(ctx context.Context, r interface{}) {
		err = toPanicError(r)
	})

	return handler(srv, stream)
}

type PanicHandler func(context.Context, interface{})

var additionalHandlers []PanicHandler

func InstallPanicHandler(handler PanicHandler) {
	additionalHandlers = append(additionalHandlers, handler)
}

func handleCrash(ctx context.Context, handler PanicHandler) {
	if r := recover(); r != nil {
		handler(ctx, r)

		if additionalHandlers != nil {
			for _, fn := range additionalHandlers {
				fn(ctx, r)
			}
		}
	}
}
