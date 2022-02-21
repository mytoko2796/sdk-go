package exporter

import (
	"net/http"
	"time"

	ste "github.com/mytoko2796/sdk-go/stdlib/telemetry/exporter/stats"
	"go.opencensus.io/stats/view"
)

type statsMode string

const (
	StatPrometheus  statsMode = `PROMETHEUS`
	StatDatadog     statsMode = `DATADOG`
	StatOpencensus  statsMode = `OPENCENSUS`
	StatStackdriver statsMode = `STACKDRIVER`
	StatZpage       statsMode = `ZPAGE`
)

type StatExporter interface {
	Export(mux *http.ServeMux) error
	Stop() error
}

type StatsOptions struct {
	Address           string
	Port              int
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration

	ViewReportingPeriod time.Duration
	RecordPeriod        time.Duration

	Zpage       ste.ZpageOptions
	Prometheus  ste.PrometheusOptions
	Opencensus  ste.OpencensusOptions
	Datadog     ste.DatadogOptions
	Stackdriver ste.StackdriverOptions
}

func (e *exporter) initStatsExporter(mode statsMode) StatExporter {
	switch mode {
	case StatZpage:
		return &ste.Zpage{Opt: e.opt.Stats.Zpage}
	case StatPrometheus:
		e.opt.Stats.Prometheus.OnError = e.onError()
		return &ste.Prometheus{Opt: e.opt.Stats.Prometheus}
	case StatOpencensus:
		return &ste.Opencensus{Opt: e.opt.Stats.Opencensus}
	case StatDatadog:
		e.opt.Stats.Datadog.OnError = e.onError()
		return &ste.Datadog{Opt: e.opt.Stats.Datadog}
	case StatStackdriver:
		e.opt.Stats.Stackdriver.OnError = e.onError()
		return &ste.Stackdriver{Opt: e.opt.Stats.Stackdriver}
	}
	return nil
}

func (e *exporter) initStatsExporters() []StatExporter {
	var exp []StatExporter
	if e.opt.Stats.Datadog.Enabled {
		exp = append(exp, e.initStatsExporter(StatDatadog))
	}
	if e.opt.Stats.Prometheus.Enabled {
		exp = append(exp, e.initStatsExporter(StatPrometheus))
	}
	if e.opt.Stats.Zpage.Enabled {
		exp = append(exp, e.initStatsExporter(StatZpage))
	}
	if e.opt.Stats.Opencensus.Enabled {
		exp = append(exp, e.initStatsExporter(StatOpencensus))
	}
	if e.opt.Stats.Stackdriver.Enabled {
		exp = append(exp, e.initStatsExporter(StatStackdriver))
	}
	if len(exp) > 0 {
		e.setStatsConfig()
	}
	return exp
}

func (e *exporter) setStatsConfig() {
	view.SetReportingPeriod(e.opt.Stats.ViewReportingPeriod)
}
