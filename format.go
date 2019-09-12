package golog

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

const (
	DEFAULT_FORMAT_TEMPLATE = "[%(asctime)s][%(levelname)s] %(message)s"
	ASCTIME_MARK            = "%(asctime)s"
	LEVELNAME_MARK          = "%(levelname)s"
	MSG_MARK                = "%(message)s"
	DEFAULT_TIME_LAYOUT     = "2006-01-02 15:04:05.999"
)

type Formatter interface {
	Format(b *bytes.Buffer, level string, format string, args ...interface{}) []byte
}

type DefaultFormatter struct {
	Template   string
	TimeLayout string
}

func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{Template: DEFAULT_FORMAT_TEMPLATE, TimeLayout: DEFAULT_TIME_LAYOUT}
}

func (this *DefaultFormatter) SetTemplate(template string) *DefaultFormatter {
	this.Template = template
	return this
}

func (this *DefaultFormatter) SetTimeLayout(layout string) *DefaultFormatter {
	this.TimeLayout = layout
	return this
}

func (this *DefaultFormatter) Format(b *bytes.Buffer, level string, format string, args ...interface{}) []byte {
	msg := strings.Replace(this.Template, ASCTIME_MARK, this.FormatTime(), -1)
	msg = strings.Replace(msg, LEVELNAME_MARK, level, -1)
	msg = strings.Replace(msg, MSG_MARK, fmt.Sprintf(format, args...), -1)

	b.WriteString(msg)
	b.WriteByte('\n')
	return b.Bytes()
}

func (this *DefaultFormatter) FormatTime() string {
	return time.Now().Format(this.TimeLayout)
}
