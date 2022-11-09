package tracing

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	jaegerClientConfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	jexpvar "github.com/uber/jaeger-lib/metrics/expvar"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"gitlab.com/dipper-iot/shared/logger"
	"gitlab.com/dipper-iot/shared/service"
	googleGRPC "google.golang.org/grpc"
	"os"
	"strings"
	"sync"
)

var (
	once = sync.Once{}
)

func JaegerTracer(configs ...jaegerClientConfig.Option) service.Option {
	enable := os.Getenv("JAEGER_ENABLE")
	if strings.ToUpper(enable) != "TRUE" {
		return nil
	}
	statusEnable = true
	return func(o *service.Options) {
		once.Do(func() {
			var metricsFactory metrics.Factory

			metricsBackend, success := os.LookupEnv("JAEGER_METRIC_BACKEND")
			if !success {
				metricsBackend = "expvar"
			}
			switch metricsBackend {
			case "expvar":
				metricsFactory = jexpvar.NewFactory(10) // 10 buckets for histograms
				logger.Info("Using expvar as metrics backend")
			case "prometheus":
				metricsFactory = jprom.New().Namespace(metrics.NSOptions{Name: "sayang", Tags: nil})
				logger.Info("Using Prometheus as metrics backend")
			default:
				logger.Fatal("unsupported metrics backend " + metricsBackend)
			}

			tracer, closer := NewTracer(o.Name, metricsFactory, configs...)

			o.BeforeStop(func(o *service.Options) error {
				return closer.Close()
			})

			o.OptionServerUnaryServerInterceptor(otgrpc.OpenTracingServerInterceptor(tracer))
			o.OptionServerStreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(tracer))

			o.AddOptionGrpcClient(
				googleGRPC.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
				googleGRPC.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
			)
		})
	}
}
