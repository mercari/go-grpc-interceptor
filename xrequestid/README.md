# xrequestid

xrequestid is an grpc interceptor which receives request id from metadata and set the request id to context. If request is is not found in metadata, generate a random request id by `github.com/renstrom/shortuuid`.

## El Toro Fork Changelog

- Use `uuid` instead of `shortuuid` since we are not passing it through context and not the url
- Fixed the Compilation errors due to the depricated `metadata.FromContext` function
- Added a `go.mod` and replaced imports with the go mod equivilent

## Usage

```golang
import (
 "github.com/eltorocorp/go-grpc-request-id-interceptor/xrequestid"
 "golang.org/x/net/context"
)

func main() {
 uIntOpt := grpc.UnaryInterceptor(xrequestid.UnaryServerInterceptor())
 sIntOpt := grpc.StreamInterceptor(xrequestid.StreamServerInterceptor())
 grpc.NewServer(uIntOpt, sIntOpt)
}

func foo(ctx context.Context) {
 requestID := xrequestid.FromContext(ctx)
 fmt.printf("requestID :%s", requestID)
}
```

### Chaining Request ID

If request id is passed by metadata, the request id is used as is by default. `xrequestid.ChainRequestID()` is an option to chain multiple request ids by generating new request id for each request and concatenating it to original request ids.

```golang
func main() {
 uInt := xrequestid.UnaryServerInterceptor(xrequestid.ChainRequestID()))
 sInt := xrequestid.StreamServerInterceptor(xrequestid.ChainRequestID()))
}
```

### Validating Request ID

It is important to validate request id in order to protect from abusing `X-Request-ID`. You can define own validator for request id, and set by `xrequestid.RequestIDValidator()`.

```golang
func customRequestIDValidator(requestID string) bool {
 if len(requestID) < 4 {
  return false
 }
 return true
}

func main() {
 uInt := xrequestid.UnaryServerInterceptor(
  xrequestid.RequestIDValidator(customRequestIDValidator),
 ))
 sInt := xrequestid.StreamServerInterceptor(
  xrequestid.RequestIDValidator(customRequestIDValidator),
    ))
}
```
