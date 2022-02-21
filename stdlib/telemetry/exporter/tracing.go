package exporter

import (
	l "log"

	tre "github.com/mytoko2796/sdk-go/stdlib/telemetry/exporter/tracing"
	"go.opencensus.io/trace"
)

type tracingMode string

const (
	TraceJaeger      tracingMode = `JAEGER`
	TraceDatadog     tracingMode = `DATADOG`
	TraceOpencensus  tracingMode = `OPENCENSUS`
	TraceStackdriver tracingMode = `STACKDRIVER`
	TraceZipkin      tracingMode = `ZIPKIN`
)

type TracingExporter interface {
	Export() error
	Stop() error
}

type TracingOptions struct {
	SamplingProbabilty float64

	Datadog     tre.DatadogOptions
	Zipkin      tre.ZipkinOptions
	Opencensus  tre.OpencensusOptions
	Jaeger      tre.JaegerOptions
	Stackdriver tre.StackdriverOptions
}

func (e *exporter) initTracingExporter(mode tracingMode) TracingExporter {
	switch mode {
	case TraceZipkin:
		return &tre.Zipkin{Opt: e.opt.Tracing.Zipkin, LogWriter: l.New(e.logger.PipeWriter(), "", l.LstdFlags)}
	case TraceDatadog:
		e.opt.Tracing.Datadog.OnError = e.onError()
		return &tre.Datadog{Opt: e.opt.Tracing.Datadog}
	case TraceOpencensus:
		return &tre.Opencensus{Opt: e.opt.Tracing.Opencensus}
	case TraceJaeger:
		e.opt.Tracing.Jaeger.OnError = e.onError()
		return &tre.Jaeger{Opt: e.opt.Tracing.Jaeger}
	case TraceStackdriver:
		e.opt.Tracing.Stackdriver.OnError = e.onError()
		return &tre.Stackdriver{Opt: e.opt.Tracing.Stackdriver}
	}
	return nil
}

func (e *exporter) initTracingExporters() []TracingExporter {
	var exp []TracingExporter
	if e.opt.Tracing.Datadog.Enabled {
		exp = append(exp, e.initTracingExporter(TraceDatadog))
	}
	if e.opt.Tracing.Jaeger.Enabled {
		exp = append(exp, e.initTracingExporter(TraceJaeger))
	}
	if e.opt.Tracing.Opencensus.Enabled {
		exp = append(exp, e.initTracingExporter(TraceOpencensus))
	}
	if e.opt.Tracing.Stackdriver.Enabled {
		exp = append(exp, e.initTracingExporter(TraceStackdriver))
	}
	if e.opt.Tracing.Zipkin.Enabled {
		exp = append(exp, e.initTracingExporter(TraceZipkin))
	}
	if len(exp) > 0 {
		e.setTracingConfig()
	}
	return exp
}

func (e *exporter) setTracingConfig() {
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(e.opt.Tracing.SamplingProbabilty)})
}
