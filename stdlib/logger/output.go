package logger

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mytoko2796/sdk-go/stdlib/error"
)

const (
	OutputStdout  string = `stdout`
	OutputFile    string = `file`
	OutputDiscard string = `discard`

	outputStdout  string = `[STDOUT]`
	outputFile    string = `[FILE]`
	outputDiscard string = `[DISCARD]`
	outputUnknown string = `[UNKNOWN LOG OUTPUT]`
)

var (
	errUnknownOutput = fmt.Errorf(`Unknown log Output`)
	ErrUnknownOutput = error.Wrapf(errUnknownOutput, errLogger, FAILED)
)

func (l *logrusImpl) convertAndSetOutput() {
	switch l.opt.Output {
	case OutputDiscard:
		l.log.Info(OK, infoLogger, outputDiscard)
		l.logger.SetOutput(ioutil.Discard)
	case OutputStdout:
		l.logger.SetOutput(os.Stdout)
		l.log.Info(OK, infoLogger, outputStdout)
	case OutputFile:
		f, err := os.OpenFile(l.opt.LogOutputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			err = error.Wrapf(err, errLogger, FAILED)
			l.log.Panic(err)
		}
		l.file = f
		l.logger.SetOutput(l.file)
		l.log.Info(OK, infoLogger, outputFile, l.opt.LogOutputPath)
	default:
		l.log.Panic(ErrUnknownOutput)
	}
}