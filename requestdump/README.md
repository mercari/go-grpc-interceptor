# requestdump

requestdump is an interceptor to log request/response messages including header/metadata as json with zap.

## Usage

```golang
import (
	"github.com/mercari/go-grpc-interceptor/requestdump"
)

func main() {
	zaplogger := zap.New(zap.NewJSONEncoder())
	uIntOpt := grpc.UnaryInterceptor(requestdump.UnaryServerInterceptor(requestdump.Zap(zaplogger)))
	grpc.NewServer(uIntOpt, sIntOpt)
}
```

### options

- `requestdump.Zap`
 - zap logger for requestdump
- `requestdump.RootKey`
 - dump message key name
- `requestdump.Disable`
 - disable requestdump
