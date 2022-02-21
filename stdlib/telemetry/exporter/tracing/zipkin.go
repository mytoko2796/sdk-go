package tracingexporter

import (
	"log"
	"time"

	"contrib.go.opencensus.io/exporter/zipkin"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
)

type Zipkin struct {
	LogWriter *log.Logger
	Opt       ZipkinOptions
	exp       *zipkin.Exporter
}

type ZipkinOptions struct {
	Enabled   bool
	Namespace string
	// AgentEndpoint instructs exporter to send spans to zipkin-agent at this address.
	AgentEndpoint string
	// CollectorEndpoint is the full url to send the spans to, e.g. http://localhost:9411/api/v2/spans
	CollectorEndpoint string
	// BatchSize sets the maximum batch size, after which a collect will be
	// triggered. The default batch size is 100 traces.
	BatchSize int
	// BatchInterval sets the maximum duration we will buffer traces before
	// emitting them to the collector. The default batch interval is 1 second.
	BatchInterval time.Duration
	// MaxBacklog sets the maximum backlog size. When batch size reaches this
	// threshold, spans from the beginning of the batch will be disposed.
	MaxBacklog int
	// Timeout sets maximum timeout for http request.
	Timeout time.Duration
}

func (e *Zipkin) Export() error {
	agentEndpoint, err := openzipkin.NewEndpoint(e.Opt.Namespace, e.Opt.AgentEndpoint)
	if err != nil {
		return err
	}
	reporter := zipkinHTTP.NewReporter(
		e.Opt.CollectorEndpoint,
		zipkinHTTP.Logger(e.LogWriter),
		zipkinHTTP.BatchSize(e.Opt.BatchSize),
		zipkinHTTP.BatchInterval(e.Opt.BatchInterval),
		zipkinHTTP.MaxBacklog(e.Opt.MaxBacklog),
		zipkinHTTP.Timeout(e.Opt.Timeout),
	)
	e.exp = zipkin.NewExporter(reporter, agentEndpoint)
	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(e.exp)
	return nil
}

func (e *Zipkin) Stop() error {
	trace.UnregisterExporter(e.exp)
	return nil
}
