package statexporter

import (
	"net/http"

	"go.opencensus.io/zpages"
)

type Zpage struct {
	Opt ZpageOptions
}

type ZpageOptions struct {
	Enabled bool
	Path    string
}

func (e *Zpage) Export(mux *http.ServeMux) error {
	zpages.Handle(mux, e.Opt.Path)
	return nil
}

func (e *Zpage) Stop() error {
	// not implemented
	return nil
}
