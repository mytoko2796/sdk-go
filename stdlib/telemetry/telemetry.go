package telemetry

import (
	"net"
	"net/http"
	"time"

	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	exporter "github.com/mytoko2796/sdk-go/stdlib/telemetry/exporter"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/gauge"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/view"

	"go.opencensus.io/plugin/ochttp"

	"sync"
)

const (
	idxMetric int = iota
	idxTracing
	idxProf
)

const (
	ReadTimeout  time.Duration = 1 * time.Second
	WriteTimeout time.Duration = 1 * time.Second
)

const (
	infoTeleServe string = `Telemetry Server:`
	errTeleServe  string = `%s %s Telemetry Server`

	OK     string = "[OK]"
	FAILED string = "[FAILED]"
)

var (
	once   = sync.Once{}
	server = []string{
		idxMetric:  "[METRIC]",
		idxTracing: "[TRACING]",
		idxProf:    "[PROFILER]",
	}
)

type Telemetry interface {
	WrapMuxHandler(h http.Handler) http.Handler
	WrapLocalHandler(isHealthCheckEndpoint bool, path string, h http.Handler) http.Handler
	WrapClientTransport(base http.RoundTripper) http.RoundTripper

	Serve(mode int, ln net.Listener)
	GetServers() []*Server
	Shutdown()
}

type telemetry struct {
	logger   log.Logger
	opt      Options
	servers  []*http.Server
	exp      exporter.Exporter
	clientRT http.RoundTripper
	termSig  chan struct{}
}

type Options struct {
	Exporters exporter.Options
}

type Server struct {
	Mode int
	Addr string
}

func Init(logger log.Logger, opt Options) Telemetry {
	var t *telemetry
	once.Do(func() {
		t = &telemetry{
			logger:  logger,
			opt:     opt,
			servers: make([]*http.Server, len(server)), // 0: metrics, 1: tracing << tracing will not use any http server
			exp:     exporter.Init(logger, opt.Exporters),
			termSig: make(chan struct{}, 1),
		}
		//telemetry server
		t.initTelemetryServer()

		//registering default components
		if err := gauge.Init(opt.Exporters.Stats, t.termSig); err != nil {
			logger.Fatal(err)
		}

		if err := view.Init(); err != nil {
			logger.Fatal(err)
		}
	})

	return t
}

func (t *telemetry) WrapLocalHandler(isHealthCheckEndpoint bool, path string, h http.Handler) http.Handler {
	if !isHealthCheckEndpoint {
		return ochttp.WithRouteTag(h, path)
	}
	return h
}

func (t *telemetry) WrapMuxHandler(h http.Handler) http.Handler {
	return &ochttp.Handler{
		Handler: h,
	}
}

func (t *telemetry) WrapClientTransport(base http.RoundTripper) http.RoundTripper {
	return &ochttp.Transport{
		Base: base,
	}
}
