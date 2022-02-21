package logger

import (
	"fmt"
	lr "github.com/sirupsen/logrus"
	"github.com/mytoko2796/sdk-go/stdlib/error"
)

const (
	FormatJSON string = "json"
	FormatText string = "text"

	formatJSON          string = "[JSON]"
	formatText          string = "[TEXT]"
	formatUnknownstring        = "[UNKNOWN LOG FORMAT]"
)

var (
	errUnknownFormat = fmt.Errorf(`Unknown log format`)
	ErrUnknownFormat = error.Wrapf(errUnknownFormat, errLogger, FAILED)
)

func (l *logrusImpl) convertAndSetFormatter() {
	switch l.opt.Formatter {
	case FormatText:
		l.logger.SetFormatter(&lr.TextFormatter{})
		l.log.Info(OK, infoLogger, formatText)
	case FormatJSON:
		l.logger.SetFormatter(&lr.JSONFormatter{})
		l.log.Info(OK, infoLogger, formatJSON)
	default:
		l.log.Panic(ErrUnknownFormat)
	}
}
