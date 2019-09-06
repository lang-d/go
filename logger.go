package golog

import (
	"io"
	"os"
)

type Logger struct {
	Level            int
	Writer           io.Writer
	console          io.Writer
	errorConsole     io.Writer
	SplitLogFileType string
}

func NewLogger() *Logger {
	return &Logger{
		Level:        INFO,
		console:      os.Stdout,
		errorConsole: os.Stderr,
	}
}

func (this *Logger) SetWriter(writer io.Writer) *Logger {
	return this
}

func (this *Logger) SetLevel(level int) *Logger {
	this.Level = level
	return this
}

func (this *Logger) SetEnableConsole(enable bool) *Logger {
	return this
}

func (this *Logger) SetSplitLogFileType(splitType string) *Logger {
	switch splitType {
	case SPLIT_BY_DAY:
	case SPLIT_BY_HOUR:
	case SPLIT_BY_SIZE:
	default:
		panic("not support this type:" + splitType)
	}
	this.SplitLogFileType = splitType
	return this
}
