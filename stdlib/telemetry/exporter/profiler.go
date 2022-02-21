package exporter

import (
	"net/http"
	"runtime"
	"time"

	pte "github.com/mytoko2796/sdk-go/stdlib/telemetry/exporter/profiler"
)

type profMode string

const (
	ProfPprof profMode = `PPROF`
)

type ProfilerExporter interface {
	Export(mux *http.ServeMux) error
	Stop() error
}

type ProfilerOptions struct {
	Address              string
	Port                 int
	ReadHeaderTimeout    time.Duration
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	MutexProfileFraction int
	Pprof                pte.PProfOptions
}

func (e *exporter) initProfilerExporter(mode profMode) ProfilerExporter {
	switch mode {
	case ProfPprof:
		return &pte.Pprof{Opt: e.opt.Profiler.Pprof}
	}
	return nil
}

func (e *exporter) initProfilerExporters() []ProfilerExporter {
	var exp []ProfilerExporter
	if e.opt.Profiler.Pprof.Enabled {
		exp = append(exp, e.initProfilerExporter(ProfPprof))
	}
	if len(exp) > 0 {
		e.setProfilerConfig()
	}
	return exp
}

func (e *exporter) setProfilerConfig() {
	runtime.SetMutexProfileFraction(e.opt.Profiler.MutexProfileFraction)
}
