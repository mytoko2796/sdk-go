package statexporter

import (
	"net/http"

	datadog "github.com/DataDog/opencensus-go-exporter-datadog"
	"go.opencensus.io/stats/view"
)

type Datadog struct {
	Opt DatadogOptions
	exp *datadog.Exporter
}

type DatadogOptions struct {
	Enabled bool
	// Namespace specifies the namespaces to which metric keys are appended.
	Namespace string
	// StatsAddr specifies the host[:port] address for DogStatsD. It defaults
	// to localhost:8125.
	AgentEndpoint string
	// OnError specifies a function that will be called if an error occurs during
	// processing stats or metrics.
	OnError func(err error)
	// Tags specifies a set of global tags to attach to each metric.
	Tags []string
	// GlobalTags holds a set of tags that will automatically be applied to all
	// exported spans.
	GlobalTags map[string]interface{}
}

func (e *Datadog) Export(mux *http.ServeMux) error {
	var err error
	e.exp, err = datadog.NewExporter(datadog.Options{
		Namespace:  e.Opt.Namespace,
		Service:    e.Opt.Namespace,
		StatsAddr:  e.Opt.AgentEndpoint,
		OnError:    e.Opt.OnError,
		Tags:       e.Opt.Tags,
		GlobalTags: e.Opt.GlobalTags,
	})
	if err != nil {
		return err
	}
	// Register it as a metrics exporter
	view.RegisterExporter(e.exp)
	return nil
}

func (e *Datadog) Stop() error {
	e.exp.Stop()
	view.UnregisterExporter(e.exp)
	return nil
}
