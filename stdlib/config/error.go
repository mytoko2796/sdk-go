package config

import "github.com/mytoko2796/sdk-go/stdlib/error"

// Ecode defines package internal error code
const (
	// Config Error Codes
	EcodeBadInput = error.Code(iota)
	EcodeTimeout
	EcodeInvalidDest
	EcodeInvalidSource
)

