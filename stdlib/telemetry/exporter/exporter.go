package exporter

import (
	"net/http"

	log "github.com/mytoko2796/sdk-go/stdlib/logger"
)

type exporter struct {
	logger log.Logger
	opt    Options
	stat   []StatExporter
	trace  []TracingExporter
	prof   []ProfilerExporter
}

type Exporter interface {
	ExportAllProfilers(mux *http.ServeMux) bool
	ExportAllStats(mux *http.ServeMux) bool
	ExportAllTracers() bool
	StopAllProfilers()
	StopAllStats()
	StopAllTracers()
}

type Options struct {
	Stats    StatsOptions
	Tracing  TracingOptions
	Profiler ProfilerOptions
}

func Init(logger log.Logger, opt Options) Exporter {
	return &exporter{
		logger: logger,
		opt:    opt,
		stat:   nil,
		trace:  nil,
		prof:   nil,
	}
}
func (e *exporter) ExportAllProfilers(mux *http.ServeMux) bool {
	e.prof = e.initProfilerExporters()
	for _, ep := range e.prof {
		if err := ep.Export(mux); err != nil {
			e.logger.Panic(err)
		}
	}
	if len(e.prof) > 0 {
		return true
	}
	return false
}

func (e *exporter) ExportAllStats(mux *http.ServeMux) bool {
	e.stat = e.initStatsExporters()
	for _, es := range e.stat {
		if err := es.Export(mux); err != nil {
			e.logger.Panic(err)
		}
	}
	if len(e.stat) > 0 {
		return true
	}
	return false
}

func (e *exporter) ExportAllTracers() bool {
	e.trace = e.initTracingExporters()
	for _, et := range e.trace {
		if err := et.Export(); err != nil {
			e.logger.Panic(err)
		}
	}
	if len(e.trace) > 0 {
		return true
	}
	return false
}

func (e *exporter) StopAllProfilers() {
	for _, ep := range e.prof {
		ep.Stop()
	}
}

func (e *exporter) StopAllStats() {
	for _, es := range e.stat {
		es.Stop()
	}
}

func (e *exporter) StopAllTracers() {
	for _, et := range e.trace {
		et.Stop()
	}
}

func (e *exporter) onError() func(err error) {
	return func(err error) {
		e.logger.Error(err)
	}
}
