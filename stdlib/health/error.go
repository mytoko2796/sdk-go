package health

import (
	errors "github.com/mytoko2796/sdk-go/stdlib/error"
)

// Ecode defines package internal error code
const (
	// Health Error Codes
	EcodeAppInitFailed = errors.Code(iota)
	EcodeReadyCheckFuncIsNil
	EcodeHealthyCheckFuncIsNil
	EcodeNotHealthy
	EcodeNotReady
	EcodeNotReadyAndHealthy
)

var (
	// Error Health Package
	ErrAppInitFailed        = errors.NewWithCode(EcodeAppInitFailed, `MaxWaitingTime is elapsed before App is ready & healthy! - Application cannot be started`, errHealth, _FAILED)
	ErrReadyCheckFuncIsNil  = errors.NewWithCode(EcodeReadyCheckFuncIsNil, `Readinesscheck Function is Nil`, errHealth, _FAILED)
	ErrHealthCheckFuncIsNil = errors.NewWithCode(EcodeHealthyCheckFuncIsNil, `Healthcheck Function is Nil`, errHealth, _FAILED)
	ErrNotHealthy           = errors.NewWithCode(EcodeNotHealthy, `app is not healthy`, errHealth, _FAILED)
	ErrNotReady             = errors.NewWithCode(EcodeNotReady, `app is not ready`, errHealth, _FAILED)
	ErrNotReadyAndHealthy   = errors.NewWithCode(EcodeNotReadyAndHealthy, `app is not ready nor healthy`, errHealth, _FAILED)
)