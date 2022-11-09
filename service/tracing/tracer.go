package tracing

import (
	"github.com/opentracing/opentracing-go"
	jaegerClientConfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
	"gitlab.com/dipper-iot/shared/logger"
	"io"
	"time"
)

func NewTracer(serviceName string, metricsFactory metrics.Factory, configs ...jaegerClientConfig.Option) (opentracing.Tracer, io.Closer) {
	traceCfg := &jaegerClientConfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerClientConfig.SamplerConfig{
			Type:  "const",
			Param: 1.0,
		},
		RPCMetrics: true,
	}

	_, err := traceCfg.FromEnv()
	if err != nil {
		logger.Fatal("cannot parse Jaeger env vars %s", err.Error())
	}
	// TODO(ys) a quick hack to ensure random generators get different seeds, which are based on current time.
	time.Sleep(100 * time.Millisecond)

	if metricsFactory != nil {
		metricsFactory = metricsFactory.Namespace(metrics.NSOptions{Name: serviceName, Tags: nil})
		configs = append(configs, jaegerClientConfig.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)))
	}

	tracer, closer, err := traceCfg.NewTracer(
		configs...,
	)
	if err != nil {
		panic("Failed to initialize tracer")
	}
	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}
