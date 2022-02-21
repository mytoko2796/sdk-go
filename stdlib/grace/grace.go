package grace

import (
	"fmt"
	"github.com/cloudflare/tableflip"
	"os"
	"os/signal"
	"syscall"
	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry"
	"github.com/mytoko2796/sdk-go/stdlib/httpserver"
	errors "github.com/mytoko2796/sdk-go/stdlib/error"
	"time"
)

const (
	infoGrace string = `Grace Upgrader:`
	errGrace  string = `%s Grace Upgrader Error`
	infoTCP   string = `TCP Listen`
	errTCP    string = `%s TCP Listen Error`

	_UPGRADE string = `[UPGRADED]`
	_OK      string = `[OK]`
	_FAILED  string = `[FAILED]`
)


type App interface {
	// Start serving app
	Serve()
	// Stop stopping app
	Stop()
}

type Options struct {
	// Pidfile define custom Pidfile, default ""
	Pidfile string
	// UpgradeTimeout wait period for new app to be ready
	UpgradeTimeout time.Duration
	// ShutdownTimeout wait period to shutdown old app
	ShutdownTimeout time.Duration
	// Network define tcp or udp. (default: tcp)
	Network string
}

type app struct {
	logger     log.Logger
	telemetry  telemetry.Telemetry
	httpServer httpserver.HTTPServer
	Upgrader   *tableflip.Upgrader
	Error      error
	Options    Options
	SigHUP     chan os.Signal
}

// Init initialize GraceApp Upgrader and starts telemetry servers and http servers.
func Init(logger log.Logger, tele telemetry.Telemetry, httpserver httpserver.HTTPServer, opt Options) App {
	upg, err := tableflip.New(tableflip.Options{UpgradeTimeout: opt.UpgradeTimeout, PIDFile: opt.Pidfile})
	if err != nil {
		err = errors.WrapWithCode(err, EcodeTableFlipFailed, errGrace, _FAILED)
		logger.Fatal(err)
	}
	logger.Info(_OK, infoGrace)

	if opt.Network == "" {
		opt.Network = `tcp`
	}
	gs := &app{
		logger:     logger,
		telemetry:  tele,
		httpServer: httpserver,
		Upgrader:   upg,
		Error:      err,
		Options:    opt,
		SigHUP:     make(chan os.Signal, 0),
	}
	signal.Notify(gs.SigHUP, syscall.SIGHUP)
	go gs.sighup()
	return gs
}

// Stop stops apps and its upgrader including all http(s) servers and configured telemetry servers if any
func (g *app) Stop() {
	g.Upgrader.Stop()
}

// sighup handle sighup signal
func (g *app) sighup() {
	for range g.SigHUP {
		err := g.Upgrader.Upgrade()
		if err != nil {
			err = errors.WrapWithCode(err, EcodeUpgradeFailed, errGrace, _FAILED)
			g.logger.Fatal(err)
			continue
		}
		g.logger.Info(_UPGRADE, infoGrace)
	}
}

// Server starts apps including all http(s) servers and configured telemetry servers if any
func (g *app) Serve() {
	// get all http(s) servers
	for _, s := range g.httpServer.GetServers() {
		ln, err := g.Upgrader.Fds.Listen(g.Options.Network, s.Addr)
		if err != nil {
			err = errors.WrapWithCode(err, EcodeHTTPServerFailed, errTCP, _FAILED)
			g.logger.Fatal(err)
		}
		g.logger.Info(_OK, infoTCP, fmt.Sprintf("@%s", s.Addr))
		go g.httpServer.Serve(s.Mode, ln)
	}

	//get all telemetry servers
	telemetryServers := g.telemetry.GetServers()
	for _, s := range telemetryServers {
		ln, err := g.Upgrader.Fds.Listen(g.Options.Network, s.Addr)
		if err != nil {
			err = errors.WrapWithCode(err, EcodeTelemetryServerFailed, errTCP, _FAILED)
			g.logger.Fatal(err)
		}
		g.logger.Info(_OK, infoTCP, fmt.Sprintf("@%s", s.Addr))
		go g.telemetry.Serve(s.Mode, ln)
	}

	if err := g.Upgrader.Ready(); err != nil {
		err = errors.WrapWithCode(err, EcodeAppNotReady, errGrace, _FAILED)
		g.logger.Fatal(err)
	}

	<-g.Upgrader.Exit()

	time.AfterFunc(g.Options.ShutdownTimeout, func() {
		err := errors.NewWithCode(EcodeAppShutdownFailed, `force Shutdown`, errGrace, _FAILED)
		g.logger.Fatal(err)
		os.Exit(1)
	})

	g.httpServer.Shutdown()
	g.telemetry.Shutdown()
}

