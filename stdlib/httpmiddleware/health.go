package httpmiddleware

import (
	"net/http"

	errors "github.com/mytoko2796/sdk-go/stdlib/error"
)

func (m *httpMiddleware) Healthcheck(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := m.healt.IsReadyAndHealthy(); err != nil {
			err = errors.WrapWithCode(err, EcodeHealth, errHealth)
			m.logger.ErrorWithContext(r.Context(), err)
			http.Error(w, "Service is Unavailable", http.StatusServiceUnavailable)
			return
		}
		fn(w, r)
	}
}

