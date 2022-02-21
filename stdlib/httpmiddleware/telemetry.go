package httpmiddleware

import (
	"fmt"
	"net/http"
	"sync"
	"context"
	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	"github.com/mytoko2796/sdk-go/stdlib/httpheader"
)

const (
	infoRespDump string = `http response sent: `
)

type telemetryResponseWriter struct {
	once       *sync.Once
	ctx        context.Context
	writer     http.ResponseWriter
	logger     log.Logger
	statusCode int
	skipLog    bool
}

func (r *telemetryResponseWriter) Header() http.Header {
	return r.writer.Header()
}
func (r *telemetryResponseWriter) Write(data []byte) (int, error) {
	n, err := r.writer.Write(data)
	if err != nil {
		return 0, err
	}
	return r.AfterWrite(data, n)
}

func (r *telemetryResponseWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode

	reqID := r.ctx.Value(httpheader.RequestID)
	if reqID != nil {
		r.writer.Header().Add(httpheader.RequestID, reqID.(string))
	}
	r.writer.WriteHeader(statusCode)
}

func (r *telemetryResponseWriter) AfterWrite(data []byte, n int) (int, error) {
	if !r.skipLog {
		r.logger.InfoWithContext(r.ctx, infoRespDump, fmt.Sprintf("status=%v payload=%s", r.statusCode, data))
	}
	return n, nil
}


