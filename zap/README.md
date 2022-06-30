# zap

Under development.

## Usage

```golang
import (
	"go.uber.org/zap"
	grpczap "github.com/mercari/go-grpc-interceptor/zap"
	"golang.org/x/net/context"
)

func main() {
	logger = zap.New(zap.NewJSONEncoder())
	uIntOpt := grpc.UnaryInterceptor(grpczap.UnaryServerInterceptor(logger))
	sIntOpt := grpc.StreamInterceptor(grpczap.StreamServerInterceptor(logger))
	grpc.NewServer(uIntOpt, sIntOpt)
}
```

### Structured logging with context

```golang
import (
	"go.uber.org/zap"
	"github.com/mercari/go-grpc-interceptor/zap/zapctx"
	"golang.org/x/net/context"
)

func foo(ctx context.Context) {
	// create new a context with some new zap.Field
	newctx := zapctx.MustNewContextWith(ctx,
		zap.String("user_id", 123456"),
	)

	// call other function with the context
	bar(newctx)
}

func bar(ctx context.Context) {
	// logging with additional context
	logger := zapctx.MustFromContext(ctx)
	logger.Info("message",
		zap.String("function", "bar"),
	)
}
```
