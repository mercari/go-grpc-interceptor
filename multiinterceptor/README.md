# multiinterceptor

multiinterceptor is a simple library which provides a function to chain multiple `UnaryServerInterceptor`s or `StreamServerInterceptor`s for gRPC.

## Usage

```go
import (
	multiint "github.com/mercari/go-grpc-interceptor/multiinterceptor"
)

func main() {
	uIntOpt := grpc.UnaryInterceptor(multiint.NewMultiUnaryServerInterceptor(
		fooUnaryInterceptor,
		barUnaryInterceptor,
	))
	sIntOpt := grpc.StreamInterceptor(multiint.NewMultiStreamServerInterceptor(
		fooStreamInterceptor,
		barStreamInterceptor,
	))
	grpc.NewServer(uIntOpt, sIntOpt)
}
```

## Context

### Unary RPC

In unary RPC, `context.Context` is passed via argument from preceding interceptor or caller.
And you can simply wrap context with `context.WithValue` and pass it to next handler.

```go
func ExampleUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	newctx := context.WithValue(ctx, "some_key", "some_value")
	return handler(newctx, req)
}
```


### Streaming RPC

In streaming RPC, `context.Context` is handled by `grpc.ServerStream` passed via argument. Call `stream.Context()` to get `context.Context` from stream.
If you want to create a new context from the context, use `NewServerStreamWithContext`, which wraps the `grpc.ServerStream` with a new context. Then you pass the wrapped stream to next handler.

```go
func ExampleStreamingInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := stream.Context()
	newctx := context.WithValue(ctx, "some_key", "some_value")
	newStream := multiint.NewServerStreamWithContext(stream, newctx)
	return handler(srv, newStream)
}
```
