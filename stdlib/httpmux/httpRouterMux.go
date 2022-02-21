package httpmux

import (
	"fmt"

	swagger "github.com/swaggo/http-swagger"

	"net/http"
	"text/template"
)

func (m *httpRouterMux) Handler() http.Handler {
	if m.cors != nil {
		return m.cors.Handler(m.mux)
	}
	return m.mux
}

func (m *httpRouterMux) handleFunc(isHealthCheckEndpoint bool, method Method, path string, hf http.HandlerFunc) {
	if !isHealthCheckEndpoint {
		hf = m.mw.WrapWithRequestContext(string(method), path, hf)
	}

	//return http.Handler wrap it using ochttp.Handler
	h := m.tele.WrapLocalHandler(isHealthCheckEndpoint, path, hf)
	m.mux.Add(string(method), path, h)
}

func (m *httpRouterMux) HandleFunc(method Method, path string, hf http.HandlerFunc) {
	m.handleFunc(false, method, path, hf)
}

func (m *httpRouterMux) registerHTTPProbeHandler() {
	m.logger.Info(OK, infoMux, fmt.Sprintf("@%s @%s", m.health.ReadyEndpoint(), m.health.HealthEndpoint()))
	m.handleFunc(true, GET, m.health.ReadyEndpoint(), m.Ready)
	m.handleFunc(true, GET, m.health.HealthEndpoint(), m.Health)
}

func (m *httpRouterMux) registerHTTPPlatformInfo() {
	if m.opt.Platform.Enabled {
		m.handleFunc(true, GET, m.opt.Platform.Path, m.conf.HTTPHandler())
		//m.handleFunc(true, GET, m.opt.Platform.PathRemote, m.remoteConf.HTTPHandler())
	}
}

func (m *httpRouterMux) registerHTTPSwagger() {
	if m.opt.Swagger.Enabled {
		m.mux.Get(m.opt.Swagger.Path, swagger.Handler(
			swagger.URL(m.opt.Swagger.DocFile),
		))
		if m.opt.Swagger.SwaggerTemplate.Enabled {
			m.handleFunc(true, GET, m.opt.Swagger.SwaggerTemplate.Path, m.swaggerTemplate)
		}
	}
}

func (m *httpRouterMux) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(healthHTTPStatus(m.health.IsHealthy()))
}

func (m *httpRouterMux) Ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(readyHTTPStatus(m.health.IsReady()))
}

func (m *httpRouterMux) swaggerTemplate(w http.ResponseWriter, r *http.Request) {
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

