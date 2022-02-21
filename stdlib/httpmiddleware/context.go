package httpmiddleware

import (
	"context"
	"fmt"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"github.com/mytoko2796/sdk-go/stdlib/httpheader"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"go.opencensus.io/trace/propagation"
)

var defaultFormat propagation.HTTPFormat = &b3.HTTPFormat{}

const (
	AppendContext    string = `Debug Append Context %s`
	MethodNotAllowed string = `Method Not Allowed`

	HTTP  string = `HTTP`
	HTTPS string = `HTTPS`
)

func (m *httpMiddleware) AppendRequestContext(method string, keyServerRoute string, fn http.HandlerFunc) MiddlewareHandle {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if method != r.Method {
				http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
				return
			}

			scheme := HTTP
			if r.TLS != nil {
				scheme = HTTPS
			}

			reqID := r.Header.Get(httpheader.RequestID)
			if reqID == "" {
				reqID = uuid.New().String()
			}

			addrIp := r.Header.Get(httpheader.ForwardedFor)
			if addrIp == "" {
				addrIp = r.RemoteAddr
			}

			ctx := context.WithValue(r.Context(), httpheader.RequestMethod, r.Method)
			ctx = context.WithValue(ctx, httpheader.RequestScheme, scheme)
			ctx = context.WithValue(ctx, httpheader.KeyServerRoute, keyServerRoute)
			ctx = context.WithValue(ctx, httpheader.ForwardedFor, addrIp)
			ctx = context.WithValue(ctx, httpheader.RequestID, reqID)

			spanContext, ok := defaultFormat.SpanContextFromRequest(r)
			if ok {
				ctx = context.WithValue(ctx, b3.TraceIDHeader, spanContext.TraceID)
			}

			m.logger.DebugWithContext(ctx, AppendContext, fmt.Sprintf("%s IpAddr=%s Method=%s Scheme=%s KeyServerRoute=%s ReqID=%s", OK, addrIp, method, scheme, keyServerRoute, reqID))

			var skip bool
			if _, ok := m.responsePathBlackList[r.Method+":"+keyServerRoute]; ok {
				skip = true
			}

			trw := &telemetryResponseWriter{
				once:    &sync.Once{},
				ctx:     ctx,
				writer:  w,
				logger:  m.logger,
				skipLog: skip,
			}

			fn(trw, r.WithContext(ctx))
		}
	}
}
