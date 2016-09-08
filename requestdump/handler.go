package requestdump

import (
	"github.com/mercari/go-grpc-interceptor/zap/zapctx"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(opt ...Option) grpc.UnaryServerInterceptor {
	opts := newOptions()
	for _, o := range opt {
		o.apply(&opts)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if !opts.enabled {
			return handler(ctx, req)
		}

		logger := opts.logger
		if logger == nil {
			logger = zapctx.MustFromContext(ctx)
		}
		dump(ctx, opts, logger, info, true, req, nil)
		resp, err = handler(ctx, req)
		dump(ctx, opts, logger, info, false, resp, err)
		return
	}
}
