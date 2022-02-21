package logger

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"sync"

	lr "github.com/sirupsen/logrus"
)

const (
	infoLogger string = `Logger:`
	errLogger  string = `%s Logger Error`
	OK         string = "[OK]"
	FAILED     string = "[FAILED]"
)

var (
	once = sync.Once{}
)


type Logger interface {
	SetOptions(opt Options)
	Stop()
	PipeWriter() io.Writer
	TraceWithContext(ctx context.Context, v ...interface{})
	DebugWithContext(ctx context.Context, v ...interface{})
	InfoWithContext(ctx context.Context, v ...interface{})
	WarnWithContext(ctx context.Context, v ...interface{})
	ErrorWithContext(ctx context.Context, v ...interface{})
	FatalWithContext(ctx context.Context, v ...interface{})
	PanicWithContext(ctx context.Context, v ...interface{})
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Panic(v ...interface{})
}

type logrusImpl struct {
	mu     *sync.RWMutex
	logger *lr.Logger
	log    *lr.Entry
	opt    Options
	file   *os.File
	writer *io.PipeWriter
}

type Options struct {
	Output string
	Formatter string
	Level string
	LogOutputPath string
	DefaultFields map[string]string
	ContextFields map[string]string
}

func Init(opt Options) Logger{
	var lg *logrusImpl
	once.Do(func() {
		logrus := lr.New()
		log := logrus.WithFields(lr.Fields{})
		lg = &logrusImpl{
			mu:     &sync.RWMutex{},
			logger: logrus,
			log:    log,
			opt:    opt,
		}

		lg.logger.SetOutput(ioutil.Discard)
		lg.setDefaultOptions()
		lg.applyOptions()
	})

	return lg
}

func (l *logrusImpl) SetOptions(opt Options) {
	l.mu.Lock()
	l.opt = opt
	l.mu.Unlock()
	l.applyOptions()
}

func (l *logrusImpl) applyOptions() {
	l.convertAndSetOutput()
	l.convertAndSetFormatter()
	l.convertAndSetLevel()
	l.setDefaultFields()
}

func (l *logrusImpl) setDefaultOptions() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.opt.Output == "" {
		//never put default to discard, error will not be displayed!
		l.opt.Output = OutputStdout
	}
	if l.opt.Formatter == "" {
		l.opt.Formatter = FormatText
	}
	if l.opt.Level == "" {
		l.opt.Level = LevelTrace
	}
}
func (l *logrusImpl) Stop() {
	l.writer.Close()
	if l.file != nil {
		l.file.Close()
	}
}

