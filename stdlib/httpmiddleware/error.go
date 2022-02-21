package httpmiddleware

import errors "github.com/mytoko2796/sdk-go/stdlib/error"

const (
	EcodeSecure = errors.Code(iota)
	EcodeHealth
	EcodePanic
	EcodeRequestDump
)

const (
	errPanic  string = `HTTP handler panic: %s full trace :%s`
	errDump   string = `HTTP request dump error`
	errHealth string = `Health Middleware Error`
)

