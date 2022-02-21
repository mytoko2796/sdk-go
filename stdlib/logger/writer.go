package logger

import (
	lr "github.com/sirupsen/logrus"
	"io"
)

func (l *logrusImpl) PipeWriter() io.Writer {
	l.writer = l.log.WriterLevel(lr.WarnLevel)
	return l.writer
}

