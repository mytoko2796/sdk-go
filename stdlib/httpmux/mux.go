package httpmux

import (
	"github.com/bmizerany/pat"
	"github.com/mytoko2796/sdk-go/stdlib/config"
	"github.com/mytoko2796/sdk-go/stdlib/health"
	"github.com/mytoko2796/sdk-go/stdlib/httpmiddleware"
	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry"
	"github.com/rs/cors"
	"net/http"
	"sync"
)

type Method string
type Mode int

const (
	// HTTPROUTER implements http pat router.
	// Supports url param composition. (e.g. /user/:userid).
	HTTPROUTER Mode = iota
	// HTTPMUXSTD implements http standard mux.
	// It does not support url param composition.
	HTTPMUXSTD
)

// HTTPMethod
const (
	POST    Method = http.MethodPost
	GET     Method = http.MethodGet
	PUT     Method = http.MethodPut
	DELETE  Method = http.MethodDelete
	OPTIONS Method = http.MethodOptions
	CONNECT Method = http.MethodConnect
	HEAD    Method = http.MethodHead
	PATCH   Method = http.MethodPatch
)

const (
	infoMux string = `Mux:`

	OK     string = "[OK]"
	FAILED string = "[FAILED]"
)


var (
	once = &sync.Once{}
	muxs = map[Mode]string{
		HTTPMUXSTD: "[HTTPMUXSTD]",
		HTTPROUTER: "[HTTPROUTER]",
	}
)

// HttpMux
type HttpMux interface {
	// HTTPROUTER implements http pat router.
	// Supports url param composition. (e.g. /user/:userid).
	HandleFunc(method Method, path string, hf http.HandlerFunc)
	// HTTPMUXSTD implements http standard mux.
	// It does not support url param composition.
	Handler() http.Handler
}



type Options struct {
	// Platform
	Platform PlatformOptions

	// Cors
	Cors CorsOptions

	// Swagger
	Swagger SwaggerOptions
}

// PlatformOptions
type PlatformOptions struct {
	Enabled bool
	// Path defines endpoint
	Path string
	// Path Runtime Config endpoint
	PathRemote string
}

// SwaggerOptions
type SwaggerOptions struct {
	Enabled bool
	// Path defines endpoint for swagger
	Path string
	// DocFile defines file path location of swaggerdoc
	DocFile string
	// Custom swagger
	SwaggerTemplate SwaggerTemplateOptions
}

type SwaggerTemplateOptions struct {
	Enabled bool
	TemplateFile string
	Path string
	GoTemplate GoTemplateOptions
}

type GoTemplateOptions struct {
	Description string
	Title string
	Version string
	Host string
	BasePath string
	OAuth2ApplicationTokenUrl string
	OAuth2PasswordTokenUrl string
}

// CorsOptions
type CorsOptions struct {
	Enabled            bool
	Mode               string
	AllowedOrigins     []string
	AllowedMethods     []string
	AllowedHeaders     []string
	ExposedHeaders     []string
	MaxAge             int
	AllowCredentials   bool
	OptionsPassthrough bool
	Debug              bool
}

type httpMux struct {
	logger     log.Logger
	mw         httpmiddleware.HttpMiddleware
	mux        *http.ServeMux
	tele       telemetry.Telemetry
	health     health.Health
	opt        Options
	conf       config.Conf
	remoteConf config.Conf
	cors       *cors.Cors
}

type httpRouterMux struct {
	logger     log.Logger
	mw         httpmiddleware.HttpMiddleware
	mux        *pat.PatternServeMux
	tele       telemetry.Telemetry
	health     health.Health
	opt        Options
	conf       config.Conf
	remoteConf config.Conf
	cors       *cors.Cors
}

// Init returns HTTPMux object. It automatically applies all default middleware to all registered http.HandlerFunc. HTTPRouter mode should be used to handle complex case as it based on Pat HTTP Router as it supports url param composition. (e.g. /user/:userid). HTTP standard mux can be used only for basic url which does not require param compositions in url.
// It requires:
// 	Logger to log.
// 	Config which will be used to show configuration on defined endpoint in platform configuration.
//	Middleware will be used to wrap http handler with middlewares.
//	Telemetry will be used to wrap http handler.
//	Health will be used to handle healtcheck/readiness check.
func Init(mode Mode, logger log.Logger, conf config.Conf, remote config.Conf, mw httpmiddleware.HttpMiddleware, tele telemetry.Telemetry, hl health.Health, opt Options) HttpMux {
	var mux HttpMux
	once.Do(func() {
		//Add default middleware
		// panic recovery
		// request-id
		// request Logger
		// security
		mw.UseDefaultMiddleware()

		var c *cors.Cors
		if opt.Cors.Enabled {
			switch opt.Cors.Mode {
			case "custom":
				c = cors.New(cors.Options{
					AllowedOrigins:     opt.Cors.AllowedOrigins,
					AllowedMethods:     opt.Cors.AllowedMethods,
					AllowedHeaders:     opt.Cors.AllowedHeaders,
					ExposedHeaders:     opt.Cors.ExposedHeaders,
					MaxAge:             opt.Cors.MaxAge,
					AllowCredentials:   opt.Cors.AllowCredentials,
					OptionsPassthrough: opt.Cors.OptionsPassthrough,
					Debug:              opt.Cors.Debug,
				})
			case "allowall":
				c = cors.AllowAll()
			case "default":
				c = cors.Default()
			default:
				c = nil
			}
		}

		switch mode {
		case HTTPMUXSTD:
			httpMuxStd := &httpMux{
				logger:     logger,
				mw:         mw,
				mux:        http.NewServeMux(),
				tele:       tele,
				health:     hl,
				opt:        opt,
				conf:       conf,
				remoteConf: remote,
				cors:       c,
			}
			httpMuxStd.registerHTTPSwagger()
			httpMuxStd.registerHTTPProbeHandler()
			httpMuxStd.registerHTTPPlatformInfo()
			mux = httpMuxStd

		case HTTPROUTER:
			httpRouterMux := &httpRouterMux{
				logger:     logger,
				mw:         mw,
				mux:        pat.New(),
				tele:       tele,
				health:     hl,
				opt:        opt,
				conf:       conf,
				remoteConf: remote,
				cors:       c,
			}
			httpRouterMux.registerHTTPSwagger()
			httpRouterMux.registerHTTPProbeHandler()
			httpRouterMux.registerHTTPPlatformInfo()
			mux = httpRouterMux
		}
		logger.Info(OK, infoMux, muxs[mode])
	})
	return mux
}