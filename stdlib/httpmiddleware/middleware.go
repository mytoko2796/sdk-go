package httpmiddleware

import (
	"github.com/mytoko2796/sdk-go/stdlib/health"
	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	"github.com/unrolled/secure"
	"net/http"
)

const (
	OK     string = "[OK]"
	FAILED string = "[FAILED]"
)


type HttpMiddleware interface {
	// UseDefaultMiddleware implements default middlewares
	//	CatchPanicAndRecover
	//	HealthCheck
	//	RequestDump
	// 	Secure
	UseDefaultMiddleware()
	// UseMiddleware append particular middlewares to exisiting middleware collections
	Use(h ...MiddlewareHandle)
	// Wrap wrap handlerfunc with handlerfunc
	Wrap(h http.HandlerFunc) http.HandlerFunc
	// Wrap handlerfunc with request context
	// method and path to differ the context and handlerfunc
	WrapWithRequestContext(method, path string, h http.HandlerFunc) http.HandlerFunc
}

type Options struct {
	Log      LoggerOptions
	Security SecurityOptions
}

type LoggerOptions struct {
	RequestPathBlackList  map[string][]string
	ResponsePathBlackList map[string][]string
}

type SecurityOptions struct {
	AllowedHosts            []string // AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
	HostsProxyHeaders       []string // HostsProxyHeaders is a set of header keys that may hold a proxied hostname value for the request.
	SSLRedirect             bool     // If SSLRedirect is set to true, then only allow HTTPS requests. Default is false.
	SSLTemporaryRedirect    bool     // If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
	SSLHost                 string
	SSLHostFunc             interface{}       // SSLHost is the host name that is used to redirect HTTP requests to HTTPS. Default is "", which indicates to use the same host.
	SSLProxyHeaders         map[string]string // SSLProxyHeaders is set of header keys with associated values that would indicate a valid HTTPS request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
	STSSeconds              int64             // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
	STSIncludeSubdomains    bool              // If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
	STSPreload              bool              // If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
	STSForceHeader          bool              // STS header is only included when the connection is HTTPS. If you want to force it to always be added, set to true. `IsDevelopment` still overrides this. Default is false.
	FrameDeny               bool              // If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
	CustomFrameOptionsValue string            // CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option. Default is "".
	ContentTypeNosniff      bool              // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
	BrowserXssFilter        bool              // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
	CustomBrowserXssValue   string            // CustomBrowserXssValue allows the X-XSS-Protection header value to be set with a custom value. This overrides the BrowserXssFilter option. Default is "".
	ContentSecurityPolicy   string            // ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "". Passing a template string will replace `$NONCE` with a dynamic nonce value of 16 bytes for each request which can be later retrieved using the Nonce function.
	PublicKey               string            // PublicKey implements HPKP to prevent MITM attacks with forged certificates. Default is "".
	ReferrerPolicy          string            // ReferrerPolicy allows the Referrer-Policy header with the value to be set with a custom value. Default is "".
	FeaturePolicy           string            // FeaturePolicy allows the Feature-Policy header with the value to be set with a custom value. Default is "".
	ExpectCTHeader          string
	IsDevelopment           bool // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
}

type httpMiddleware struct {
	logger                log.Logger
	healt                 health.Health
	security              *secure.Secure
	mwares                []MiddlewareHandle
	requestPathBlackList  map[string]bool
	responsePathBlackList map[string]bool
}

type MiddlewareHandle func(http.HandlerFunc) http.HandlerFunc

func Init(logger log.Logger, healt health.Health, opt Options) HttpMiddleware {
	var SSLHostFunc *secure.SSLHostFunc
	if opt.Security.SSLHostFunc != nil {
		if f, ok := opt.Security.SSLHostFunc.(secure.SSLHostFunc); ok {
			SSLHostFunc = &f
		}
	}
	// blacklist path to skip log
	requestPathBlackList := make(map[string]bool)
	for p, methods := range opt.Log.RequestPathBlackList {
		for _, m := range methods {
			requestPathBlackList[m+":"+p] = true
		}
	}

	responsePathBlackList := make(map[string]bool)
	for p, methods := range opt.Log.ResponsePathBlackList {
		for _, m := range methods {
			responsePathBlackList[m+":"+p] = true
		}
	}

	return &httpMiddleware{
		logger: logger,
		healt:  healt,
		security: secure.New(secure.Options{
			AllowedHosts:            opt.Security.AllowedHosts,
			HostsProxyHeaders:       opt.Security.HostsProxyHeaders,
			SSLRedirect:             opt.Security.SSLRedirect,
			SSLTemporaryRedirect:    opt.Security.SSLTemporaryRedirect,
			SSLHost:                 opt.Security.SSLHost,
			SSLHostFunc:             SSLHostFunc,
			SSLProxyHeaders:         opt.Security.SSLProxyHeaders,
			STSSeconds:              opt.Security.STSSeconds,
			STSIncludeSubdomains:    opt.Security.STSIncludeSubdomains,
			STSPreload:              opt.Security.STSPreload,
			ForceSTSHeader:          opt.Security.STSForceHeader,
			FrameDeny:               opt.Security.FrameDeny,
			CustomFrameOptionsValue: opt.Security.CustomFrameOptionsValue,
			ContentTypeNosniff:      opt.Security.ContentTypeNosniff,
			BrowserXssFilter:        opt.Security.BrowserXssFilter,
			CustomBrowserXssValue:   opt.Security.CustomBrowserXssValue,
			ContentSecurityPolicy:   opt.Security.ContentSecurityPolicy,
			PublicKey:               opt.Security.PublicKey,
			ReferrerPolicy:          opt.Security.ReferrerPolicy,
			FeaturePolicy:           opt.Security.FeaturePolicy,
			ExpectCTHeader:          opt.Security.ExpectCTHeader,
			IsDevelopment:           opt.Security.IsDevelopment,
		}),
		mwares:                nil,
		requestPathBlackList:  requestPathBlackList,
		responsePathBlackList: responsePathBlackList,
	}
}

func (m *httpMiddleware) UseDefaultMiddleware() {
	m.Use(
		m.CatchPanicAndReport,
		m.Healthcheck,
		m.RequestDump,
		m.Secure,
	)
}

func (m *httpMiddleware) Use(h ...MiddlewareHandle) {
	m.mwares = append(m.mwares, h...)
}

func (m *httpMiddleware) Reverse(h []MiddlewareHandle) []MiddlewareHandle {
	for left, right := 0, len(h)-1; left < right; left, right = left+1, right-1 {
		h[left], h[right] = h[right], h[left]
	}
	return h
}

func (m *httpMiddleware) Wrap(h http.HandlerFunc) http.HandlerFunc {
	for _, mw := range m.Reverse(m.mwares) {
		h = mw(h)
	}
	return h
}

func (m *httpMiddleware) WrapWithRequestContext(method string, path string, h http.HandlerFunc) http.HandlerFunc {
	mws := []MiddlewareHandle{m.AppendRequestContext(method, path, h)}
	mws = append(mws, m.mwares...)
	for _, mw := range m.Reverse(mws) {
		h = mw(h)
	}
	return h
}
