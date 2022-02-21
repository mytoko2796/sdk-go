package profiler

import (
	"fmt"
	"net/http"
	pprof "net/http/pprof"
)

type Pprof struct {
	Opt PProfOptions
}

type PProfOptions struct {
	PathPrefix string
	Enabled    bool
	Cmdline    bool
	Profile    bool
	Symbol     bool
	Trace      bool
}

func (e *Pprof) Export(mux *http.ServeMux) error {
	mux.HandleFunc(fmt.Sprintf("%s/", e.Opt.PathPrefix), pprof.Index)
	if e.Opt.Cmdline {
		mux.HandleFunc(fmt.Sprintf("%s/cmdline", e.Opt.PathPrefix), pprof.Cmdline)
	}
	if e.Opt.Profile {
		mux.HandleFunc(fmt.Sprintf("%s/profile", e.Opt.PathPrefix), pprof.Profile)
	}
	if e.Opt.Symbol {
		mux.HandleFunc(fmt.Sprintf("%s/symbol", e.Opt.PathPrefix), pprof.Symbol)
	}
	if e.Opt.Trace {
		mux.HandleFunc(fmt.Sprintf("%s/trace", e.Opt.PathPrefix), pprof.Trace)
	}
	return nil
}

func (e *Pprof) Stop() error {
	return nil
}
