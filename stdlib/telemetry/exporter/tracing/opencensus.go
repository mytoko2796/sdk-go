package tracingexporter

import (
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

type Opencensus struct {
	Opt OpencensusOptions
	exp *ocagent.Exporter
}

type OpencensusOptions struct {
	Enabled            bool
	Namespace          string
	Insecure           bool
	ReconnectionPeriod time.Duration
	AgentEndpoint      string
	Headers            map[string]string
	Compressor         string
}

func (e *Opencensus) Export() error {
	var (
		err        error
		ocagentOpt = []ocagent.ExporterOption{
			ocagent.WithServiceName(e.Opt.Namespace),
		}
	)

	if e.Opt.Insecure {
		ocagentOpt = append(ocagentOpt, ocagent.WithInsecure())
	}
	if e.Opt.AgentEndpoint != "" {
		ocagentOpt = append(ocagentOpt, ocagent.WithAddress(e.Opt.AgentEndpoint))
	}
	if e.Opt.ReconnectionPeriod > 0 {
		ocagentOpt = append(ocagentOpt, ocagent.WithReconnectionPeriod(e.Opt.ReconnectionPeriod))
	}
	if e.Opt.Headers != nil {
		ocagentOpt = append(ocagentOpt, ocagent.WithHeaders(e.Opt.Headers))
	}
	if e.Opt.Compressor != "" {
		ocagentOpt = append(ocagentOpt, ocagent.UseCompressor(e.Opt.Compressor))
	}

	e.exp, err = ocagent.NewExporter(ocagentOpt...)
	if err != nil {
		return err
	}
	trace.RegisterExporter(e.exp)
	return nil
}

func (e *Opencensus) Stop() error {
	e.exp.Flush()
	if err := e.exp.Stop(); err != nil {
		return err
	}
	view.UnregisterExporter(e.exp)
	return nil
}
