# instrument

under development.

## Usage

```golang
import (
	"github.com/mercari/go-grpc-interceptor/instrument"
)

func main() {
	uIntOpt := grpc.UnaryInterceptor(instrument.UnaryServerInterceptor)
	sIntOpt := grpc.StreamInterceptor(instrument.StreamServerInterceptor)
	grpc.NewServer(uIntOpt, sIntOpt)
}
```

## Built-in instrumentation

### Prometheus

```golang
import (
	"github.com/mercari/go-grpc-interceptor/instrument/prometheus"
)
```

#### metrics

- grpc_calls_total (counter)
 - The total number of gRPC calls 
- grpc_calls_errors  (counter)
 - The total number of gRPC calls that returned error
- grpc_calls_duration (summary)
 - Duration of gRPC calls
 - **optional**: Call `EnableDurationSummary` To enable this metric.
