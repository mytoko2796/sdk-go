package httpserver

import (
	errors "github.com/mytoko2796/sdk-go/stdlib/error"
)

const (
	EcodeAppShutdownFailed = errors.Code(iota)
	EcodeListenerFailed
)
const (
	errServe          string = `%s %s Server`
	errNetListener           = `failed to get %s net listener`
	errListenerFailed string = `Cannot start listening %s @%s`
)
