# panichandler

panichandler is an interceptor to protect a process from aborting by panic and return Internal error as status code.

## Usage

```
import (
 "github.com/eltorocorp/go-grpc-request-id-interceptor/panichandler"
)

func main() {
 uIntOpt := grpc.UnaryInterceptor(panichandler.UnaryServerInterceptor)
 sIntOpt := grpc.StreamInterceptor(panichandler.StreamServerInterceptor)
 grpc.NewServer(uIntOpt, sIntOpt)
}
```

## Custom Panic Handler

You can write custom panic handler in case of panic. Use `InstallPanicHandler`.

```
func main() {
 panichandler.InstallPanicHandler(func(ctx context.Context, r interface{}) {
  fmt.Printf("panic happened: %v", r)
 }
}
```

### Built-in custom panic handler

- panichandler.LogPanicDump
- `debug.Stack()` to stderr
- glog.LogPanicStackMultiLine
  - show stack trace in multi line by glog
- zap.LogPanicWithStackTrace
  - use zap.Logger in context and log panic with stack trace
