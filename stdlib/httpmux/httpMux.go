package httpmux

import (
	"fmt"
	"net/http"
	"text/template"

	swagger "github.com/swaggo/http-swagger"
)

func (m *httpMux) Handler() http.Handler {
	if m.cors != nil {
		return m.cors.Handler(m.mux)
	}
	return m.mux
}

func (m *httpMux) HandleFunc(method Method, path string, hf http.HandlerFunc) {
	m.handleFunc(false, method, path, hf)
}

func (m *httpMux) registerHTTPSwagger() {
	if m.opt.Swagger.Enabled {
		m.mux.Handle(m.opt.Swagger.Path, swagger.Handler(
			swagger.URL(m.opt.Swagger.DocFile),
		))
		if m.opt.Swagger.SwaggerTemplate.Enabled {
			m.handleFunc(true, GET, m.opt.Swagger.SwaggerTemplate.Path, m.swaggerTemplate)
		}
	}
}

func (m *httpMux) handleFunc(isHealthCheckEndpoint bool, method Method, path string, hf http.HandlerFunc) {
	if !isHealthCheckEndpoint {
		hf = m.mw.WrapWithRequestContext(string(method), path, hf)
	}
	//return http.Handler wrap it using ochttp.Handler
	//skip if it is healthcheck endpoint
	h := m.tele.WrapLocalHandler(isHealthCheckEndpoint, path, hf)
	m.mux.Handle(path, h)
}

func (m *httpMux) swaggerTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(m.opt.Swagger.SwaggerTemplate.TemplateFile)
	if err != nil {
		m.logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.Execute(w, m.opt.Swagger.SwaggerTemplate.GoTemplate); err != nil {
		m.logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func (m *httpMux) registerHTTPProbeHandler() {
	m.logger.Info(OK, infoMux, fmt.Sprintf("@%s @%s", m.health.ReadyEndpoint(), m.health.HealthEndpoint()))
	m.handleFunc(true, GET, m.health.ReadyEndpoint(), m.Ready)
	m.handleFunc(true, GET, m.health.HealthEndpoint(), m.Health)
}

func (m *httpMux) Ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(readyHTTPStatus(m.health.IsReady()))
}

func (m *httpMux) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(healthHTTPStatus(m.health.IsHealthy()))
}

func (m *httpMux) registerHTTPPlatformInfo() {
	if m.opt.Platform.Enabled {
		m.handleFunc(true, GET, m.opt.Platform.Path, m.conf.HTTPHandler())
		m.handleFunc(true, GET, m.opt.Platform.PathRemote, m.remoteConf.HTTPHandler())
	}
}
