package statexporter

import (
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
)

type PrometheusOptions struct {
	Enabled   bool
	Namespace string
	Path      string
	OnError   func(err error)
}

type Prometheus struct {
	Opt PrometheusOptions
	exp *prometheus.Exporter
}

func (e *Prometheus) Export(mux *http.ServeMux) error {
	var err error
	e.exp, err = prometheus.NewExporter(prometheus.Options{
		Namespace: e.Opt.Namespace,
		OnError:   e.Opt.OnError,
	})
	if err != nil {
		return err
	}
	mux.Handle(e.Opt.Path, e.exp)
	//register prometheus endpoint
	view.RegisterExporter(e.exp)
	return nil
}

func (e *Prometheus) Stop() error {
	//not flush method implemented
	view.UnregisterExporter(e.exp)
	return nil
}
