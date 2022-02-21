package logger

func (l *logrusImpl) setDefaultFields() {
	l.mu.RLock()
	for k, v := range l.opt.DefaultFields {
		l.log = l.log.WithField(k, v)
	}
	l.mu.RUnlock()
}

