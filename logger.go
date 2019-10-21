package golog

import "os"

type Logger struct {
	Handlers []Handler
	Level    *level
}

func NewLogger() *Logger {
	return &Logger{Handlers: []Handler{NewConsleHandler()}, Level: INFO}
}

func (this *Logger) AddHandler(handler Handler) *Logger {
	this.Handlers = append(this.Handlers, handler)
	return this
}

func (this *Logger) SetHandlers(handlers []Handler) *Logger {
	this.Handlers = handlers
	return this
}

func (this *Logger) SetLevel(lev *level) *Logger {
	this.Level = lev
	return this
}

func (this *Logger) Infof(format string, args ...interface{}) *Logger {
	return this.Log(INFO, format, args...)
}

func (this *Logger) Debugf(format string, args ...interface{}) *Logger {
	return this.Log(DEBUG, format, args...)
}

func (this *Logger) Errorf(format string, args ...interface{}) *Logger {
	return this.Log(ERROR, format, args...)
}

func (this *Logger) Exceptionf(format string, args ...interface{}) *Logger {
	return this.Log(EXCEPTION, format, args...)
}

func (this *Logger) Warnf(format string, args ...interface{}) *Logger {
	return this.Log(WARN, format, args...)
}

func (this *Logger) Fatalf(format string, args ...interface{}) {
	this.Log(FATAL, format, args...)
	os.Exit(1)
}

func (this *Logger) Log(level *level, format string, args ...interface{}) *Logger {
	if this.Level.isEnable(level) {
		for _, h := range this.Handlers {
			h.Work(level, format, args...)
		}
	}
	return this
}
