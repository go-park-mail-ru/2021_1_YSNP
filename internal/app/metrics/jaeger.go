package metrics

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

type Jaeger struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func NewJaeger(serverName string) (Jaeger, error) {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: serverName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	return Jaeger{
		tracer: tracer,
		closer: closer,
	}, err
}

func (j *Jaeger) GetTracer() opentracing.Tracer {
	return j.tracer
}

func (j *Jaeger) SetGlobalTracer() {
	opentracing.SetGlobalTracer(j.tracer)
}

func (j *Jaeger) Close() {
	j.closer.Close()
}
