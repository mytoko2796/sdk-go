package httpmiddleware

import (
	"fmt"
	"net/http"
	errors "github.com/mytoko2796/sdk-go/stdlib/error"
)

const (
	errSecure string = "Middleware Secure Error"
)

var (
	HTTPRedirect = fmt.Errorf(`redirecting to HTTPS`)
)

func (m *httpMiddleware) Secure(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := m.security.Process(w, r)
		if err != nil {
			if err != HTTPRedirect {
				err = errors.WrapWithCode(err, EcodeSecure, errSecure)
				m.logger.ErrorWithContext(r.Context(), err)
			}
			return
		}
		fn(w, r)
	}
}
