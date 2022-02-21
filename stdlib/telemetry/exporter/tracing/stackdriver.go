package tracingexporter

import (
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

type Stackdriver struct {
	Opt StackdriverOptions
	exp *stackdriver.Exporter
}

type StackdriverOptions struct {
	Enabled                  bool
	Namespace                string
	ProjectID                string
	Location                 string
	MetricPrefix             string
	DefaultTraceAttributes   map[string]interface{}
	OnError                  func(err error)
	BatchInterval            time.Duration
	Timeout                  time.Duration
	TraceSpansBufferMaxBytes int
	BundleDelayThreshold     time.Duration
	BundleCountThreshold     int
}

func (e *Stackdriver) Export() error {
	//stackdriver does not need to bind to httpMux
	var err error
	e.exp, err = stackdriver.NewExporter(stackdriver.Options{
		ProjectID:                e.Opt.ProjectID,
		Location:                 e.Opt.Location,
		MetricPrefix:             e.Opt.MetricPrefix,
		OnError:                  e.Opt.OnError,
		DefaultTraceAttributes:   e.Opt.DefaultTraceAttributes,
		ReportingInterval:        e.Opt.BatchInterval,
		Timeout:                  e.Opt.Timeout,
		BundleDelayThreshold:     e.Opt.BundleDelayThreshold,
		TraceSpansBufferMaxBytes: e.Opt.TraceSpansBufferMaxBytes,
		BundleCountThreshold:     e.Opt.BundleCountThreshold,
	})
	if err != nil {
		return err
	}
	// Register it as a metrics exporter
	trace.RegisterExporter(e.exp)
	return nil
}

func (e *Stackdriver) Stop() error {
	e.exp.Flush()
	e.exp.StopMetricsExporter()
	trace.UnregisterExporter(e.exp)
	return nil
}
