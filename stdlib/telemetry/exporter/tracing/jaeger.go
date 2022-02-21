package tracingexporter

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

type Jaeger struct {
	Opt JaegerOptions
	exp *jaeger.Exporter
}

type JaegerOptions struct {
	Enabled bool

	Namespace string
	// AgentEndpoint instructs exporter to send spans to jaeger-agent at this address.
	// For example, localhost:6831.
	AgentEndpoint string
	// CollectorEndpoint is the full url to the Jaeger HTTP Thrift collector.
	// For example, http://localhost:14268/api/traces
	CollectorEndpoint string
	// OnError is the hook to be called when there is
	// an error occurred when uploading the stats data.
	// If no custom hook is set, errors are logged.
	// Optional.
	OnError func(err error)
	// Username to be used if basic auth is required.
	// Optional.
	Username string
	// Password to be used if basic auth is required.
	// Optional.
	Password string
	//BufferMaxCount defines the total number of traces that can be buffered in memory
	BufferMaxCount int
}

func (e *Jaeger) Export() error {
	var err error
	e.exp, err = jaeger.NewExporter(jaeger.Options{
		ServiceName:       e.Opt.Namespace,
		CollectorEndpoint: e.Opt.CollectorEndpoint,
		AgentEndpoint:     e.Opt.AgentEndpoint,
		OnError:           e.Opt.OnError,
		Username:          e.Opt.Username,
		Password:          e.Opt.Password,
		BufferMaxCount:    e.Opt.BufferMaxCount,
	})
	if err != nil {
		return err
	}
	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(e.exp)
	return nil
}

func (e *Jaeger) Stop() error {
	e.exp.Flush()
	trace.UnregisterExporter(e.exp)
	return nil
}
