package httpmiddleware

import (
	"net/http"
	"runtime/debug"

	errors "github.com/mytoko2796/sdk-go/stdlib/error"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"go.opencensus.io/stats"
)

func (m *httpMiddleware) CatchPanicAndReport(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				err = errors.NewWithCode(EcodePanic, errPanic, err, debug.Stack())
				m.logger.ErrorWithContext(r.Context(), err)
				//catch total panic by http handler - path is carried within context from AppendRequestContext middleware
				stats.Record(r.Context(), stat.StatPanic.M(1))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	}
}

