package httpserver

import (
	"fmt"
	"net"
	"net/http"

	errors "github.com/mytoko2796/sdk-go/stdlib/error"
)

func (h *httpServer) GetTLSKey() string {
	return h.opt.TLSKeyFile
}

func (h *httpServer) GetTLSCert() string {
	return h.opt.TLSCertFile
}

func (h *httpServer) Serve(mode int, ln net.Listener) {
	if ln == nil {
		err := errors.Wrap(fmt.Errorf(errNetListener, server[mode]), errServe)
		h.logger.Panic(err)
	}
	h.logger.Info(OK, infoServe, fmt.Sprintf("%s @%s", server[mode], ln.Addr().String()))
	if h.servers[mode] != nil {
		if mode == HTTPS && h.opt.TLSEnabled && h.opt.TLSCertFile != "" && h.opt.TLSKeyFile != "" {
			if err := h.servers[mode].ServeTLS(ln, h.opt.TLSCertFile, h.opt.TLSKeyFile); err != http.ErrServerClosed {
				err = errors.WrapWithCode(err, EcodeListenerFailed, errServe, FAILED, server[mode])
				h.logger.Panic(err)
			}
		} else {
			if err := h.servers[mode].Serve(ln); err != http.ErrServerClosed {
				err = errors.WrapWithCode(err, EcodeListenerFailed, errServe, FAILED, server[mode])
				h.logger.Panic(err)
			}
		}
	} else {
		err := errors.NewWithCode(EcodeListenerFailed, errListenerFailed, server[mode], ln.Addr().String())
		h.logger.Panic(err)
	}
}
