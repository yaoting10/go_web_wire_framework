package job

import "goboot/pkg/log"

type Logger struct {
	Log *log.Logger
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.Log.Debugf("[JOB] - "+format, a...)
}

func (l *Logger) Error(format string, a ...interface{}) {
	l.Log.Errorf("[JOB] - "+format, a...)
}
