package httpserver

import (
	"crypto/tls"
	"fmt"
	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry"
	"github.com/mytoko2796/sdk-go/stdlib/httpmux"
	l "log"
	"net"
	"sync"
	"net/http"
	"time"
	"context"
	errors "github.com/mytoko2796/sdk-go/stdlib/error"
)

const (
	HTTP int = iota
	HTTPS
)

const (
	infoServe string = `Server:`

	OK     string = "[OK]"
	FAILED string = "[FAILED]"
)

var (
	once   = &sync.Once{}
	server = []string{
		HTTP:  "[HTTP]",
		HTTPS: "[HTTPS]",
	}
)

type HTTPServer interface {
	Serve(mode int, ln net.Listener)
	Shutdown()
	GetServers() []*Server
	GetTLSCert() string
	GetTLSKey() string
}


type Server struct {
	Mode int
	Addr string
}

type Options struct {
	Address           string
	Port              int
	TLSPort           int
	TLSEnabled        bool
	TLSCertFile       string
	TLSKeyFile        string
	TLSConfig         *tls.Config
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

type httpServer struct {
	logger  log.Logger
	servers []*http.Server
	opt     Options
}

func Init(logger log.Logger, tele telemetry.Telemetry, mux httpmux.HttpMux, opt Options) HTTPServer {
	var (
		h *httpServer
	)
	once.Do(func() {
		//init http server
		h = &httpServer{
			logger:  logger,
			opt:     opt,
			servers: nil,
		}

		//Intercept Handler with telemetry handle
		handler := tele.WrapMuxHandler(mux.Handler())

		//Init HTTP Server with Handler
		h.InitHTTPServer(handler)
	})
	return h
}

func (h *httpServer) InitHTTPServer(handler http.Handler) {
	h.servers = append(h.servers, h.NewHTTPServer(false, handler))
	if h.opt.TLSEnabled {
		h.servers = append(h.servers, h.NewHTTPServer(true, handler))
	}
}

func (h *httpServer) GetServers() []*Server {
	var servers []*Server
	for k, s := range h.servers {
		servers = append(servers, &Server{Mode: k, Addr: s.Addr})
	}
	return servers
}

func (h *httpServer) Shutdown() {
	for _, s := range h.servers {
		if err := s.Shutdown(context.Background()); err != nil {
			err = errors.WrapWithCode(err, EcodeAppShutdownFailed, errServe, FAILED, s.Addr)
			s.ErrorLog.Fatal(err)
		}
	}
}


func (h *httpServer) NewHTTPServer(tlsEnabled bool, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", h.opt.Address, h.opt.Port),
		Handler:           handler,
		ReadHeaderTimeout: h.opt.ReadHeaderTimeout,
		ReadTimeout:       h.opt.ReadTimeout * time.Second,
		WriteTimeout:      h.opt.WriteTimeout * time.Second,
		IdleTimeout:       h.opt.IdleTimeout * time.Second,
		ErrorLog:          l.New(h.logger.PipeWriter(), "", l.LstdFlags),
	}

	if tlsEnabled {
		server.Addr = fmt.Sprintf("%s:%d", h.opt.Address, h.opt.TLSPort)
		server.TLSConfig = h.opt.TLSConfig
		server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
	}
	return server
}
