package grace

import errors "github.com/mytoko2796/sdk-go/stdlib/error"

// Ecode defines package internal error code
const (
	// Grace Error Codes
	EcodeTableFlipFailed = errors.Code(iota)
	EcodeUpgradeFailed
	EcodeAppNotReady
	EcodeAppShutdownFailed
	EcodeHTTPServerFailed
	EcodeTelemetryServerFailed
)

