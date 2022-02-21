package telemetry

import (
	"context"
	"fmt"
	l "log"
	"net"
	"net/http"
	"time"

	errors "github.com/mytoko2796/sdk-go/stdlib/error"
)

const (
	errNetListener = `failed to get %s net listener`
)

func (t *telemetry) initTelemetryServer() {
	statMux := http.NewServeMux()
	if t.exp.ExportAllStats(statMux) {
		opt := t.opt.Exporters.Stats
		t.servers[idxMetric] = t.NewServer(
			statMux,
			fmt.Sprintf("%s:%d", opt.Address, opt.Port),
			opt.ReadHeaderTimeout,
			opt.ReadTimeout,
			opt.WriteTimeout)
	}

	t.exp.ExportAllTracers()
	//always nil since it requires no metrics server
	t.servers[idxTracing] = nil

	profMux := http.NewServeMux()
	if t.exp.ExportAllProfilers(profMux) {
		opt := t.opt.Exporters.Profiler
		t.servers[idxProf] = t.NewServer(
			profMux,
			fmt.Sprintf("%s:%d", opt.Address, opt.Port),
			opt.ReadHeaderTimeout,
			opt.ReadTimeout,
			opt.WriteTimeout)
	}
}

func (t *telemetry) Serve(index int, ln net.Listener) {
	if ln == nil {
		err := errors.Wrap(fmt.Errorf(errNetListener, server[index]), errTeleServe)
		t.logger.Panic(err)
	}

	t.logger.Info(OK, infoTeleServe, fmt.Sprintf("%s @%s", server[index], ln.Addr().String()))
	if t.servers[index] != nil {
		if err := t.servers[index].Serve(ln); err != http.ErrServerClosed {
			err = errors.Wrapf(err, errTeleServe, FAILED, server[index])
			t.logger.Panic(err)
		}
	} else {
		t.logger.Panic(fmt.Errorf(`Cannot start listening %s @%s`, server[index], ln.Addr().String()))
	}
}

func (t *telemetry) GetServers() []*Server {
	var servers []*Server
	for k, s := range t.servers {
		if s != nil {
			servers = append(servers, &Server{Mode: k, Addr: s.Addr})
		}
	}
	return servers
}

func (t *telemetry) Shutdown() {
	t.termSig <- struct{}{}
	t.exp.StopAllStats()
	t.exp.StopAllTracers()
	t.exp.StopAllProfilers()
	for _, s := range t.servers {
		if s != nil {
			if err := s.Shutdown(context.Background()); err != nil {
				err = errors.Wrap(err, errTeleServe)
				s.ErrorLog.Fatal(err)
			}
		}
	}
}

func (t *telemetry) NewServer(handler http.Handler, address string, readHeaderTimeout, readTimeout, writeTimeout time.Duration) *http.Server {
	return &http.Server{
		Handler:           handler,
		Addr:              address,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		ErrorLog:          l.New(t.logger.PipeWriter(), "", l.LstdFlags),
	}
}
