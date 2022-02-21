package httpmiddleware

import (
	"net/http"
	"net/http/httputil"

	errors "github.com/mytoko2796/sdk-go/stdlib/error"
	"github.com/mytoko2796/sdk-go/stdlib/httpheader"
)

const (
	infoRequestDump string = `httpserver Received Request: `
)

func (m *httpMiddleware) RequestDump(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var skip bool

		keyServerRouteCtx := r.Context().Value(httpheader.KeyServerRoute)
		if keyServerRouteCtx != nil {
			keyServerRoute, _ := keyServerRouteCtx.(string)
			if _, ok := m.requestPathBlackList[r.Method+":"+keyServerRoute]; ok {
				skip = true
			}
		}

		if r == nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !skip {
			dump, err := httputil.DumpRequest(r, true)
			if err != nil {
				err = errors.WrapWithCode(err, EcodeRequestDump, errDump)
				m.logger.ErrorWithContext(r.Context(), err)
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			m.logger.InfoWithContext(r.Context(), infoRequestDump, string(dump))
		}
		fn(w, r)
	}
}
