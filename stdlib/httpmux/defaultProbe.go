package httpmux

import (
	"net/http"
)

// health converts err returned as http header.code
func healthHTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusServiceUnavailable
}

// ready converts err returned as http header.code
func readyHTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusServiceUnavailable
}
