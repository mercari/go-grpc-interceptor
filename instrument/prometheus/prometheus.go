package prometheus

import (
	"time"

	"github.com/eltorocorp/go-grpc-request-id-interceptor/instrument"
	prom "github.com/prometheus/client_golang/prometheus"
)

func init() {
	instrument.InstallInstrumentor(promInstrument)

	prom.MustRegister(totalCalls)
	prom.MustRegister(errorCalls)
}

var (
	totalCalls = prom.NewCounterVec(prom.CounterOpts{
		Name: "grpc_calls_total",
		Help: "Number of gRPC calls.",
	}, []string{"method"})
	errorCalls = prom.NewCounterVec(prom.CounterOpts{
		Name: "grpc_calls_errors",
		Help: "Number of gRPC calls that returned error.",
	}, []string{"method"})
	durations = prom.NewSummaryVec(prom.SummaryOpts{
		Name: "grpc_calls_duration",
		Help: "Duration of gRPC calls.",
	}, []string{"method"})
)

var (
	durationSummaryEnabled = false
)

func EnableDurationSummary() {
	durationSummaryEnabled = true
	prom.MustRegister(durations)
}

func promInstrument(call instrument.MethodCall) {
	labels := prom.Labels{"method": call.FullMethod}
	totalCalls.With(labels).Inc()
	if call.Error != nil {
		errorCalls.With(labels).Inc()
	}
	if durationSummaryEnabled {
		durations.With(labels).Observe(
			float64(call.Duration.Nanoseconds() / int64(time.Millisecond)),
		)
	}
}
