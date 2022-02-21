package logger

import (
	"context"

	lr "github.com/sirupsen/logrus"
)

const (
	LevelTrace string = "trace"
	LevelDebug string = "debug"
	LevelInfo  string = "info"
	LevelWarn  string = "warn"
	LevelError string = "error"
	LevelFatal string = "fatal"
	LevelPanic string = "panic"
)

func (l *logrusImpl) parseContextFields(ctx context.Context) *lr.Entry {
	doLog := l.log
	if ctx != nil {
		for k, v := range l.opt.ContextFields {
			if val := ctx.Value(v); val != nil {
				doLog = doLog.WithField(k, val)
			}
		}
	}
	return doLog
}

func (l *logrusImpl) TraceWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Trace(v...)
}

func (l *logrusImpl) Trace(v ...interface{}) {
	l.TraceWithContext(nil, v...)
}

func (l *logrusImpl) DebugWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Debug(v...)
}

func (l *logrusImpl) Debug(v ...interface{}) {
	l.DebugWithContext(nil, v...)
}

func (l *logrusImpl) InfoWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Info(v...)
}

func (l *logrusImpl) Info(v ...interface{}) {
	l.InfoWithContext(nil, v...)
}

func (l *logrusImpl) WarnWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Warn(v...)
}

func (l *logrusImpl) Warn(v ...interface{}) {
	l.WarnWithContext(nil, v...)
}

func (l *logrusImpl) ErrorWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Error(v...)
}

func (l *logrusImpl) Error(v ...interface{}) {
	l.ErrorWithContext(nil, v...)
}

func (l *logrusImpl) FatalWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Fatal(v...)
}

func (l *logrusImpl) Fatal(v ...interface{}) {
	l.FatalWithContext(nil, v...)
}

func (l *logrusImpl) PanicWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Panic(v...)
}

func (l *logrusImpl) Panic(v ...interface{}) {
	l.PanicWithContext(nil, v...)
}
