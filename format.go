package golog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"time"
)

const (
	DEFAULT_FORMAT_TEMPLATE = "[%(asctime)s][%(levelname)s] %(message)s"
	ASCTIME_MARK            = "%(asctime)s"
	LEVELNAME_MARK          = "%(levelname)s"
	MSG_MARK                = "%(message)s"
	DEFAULT_TIME_LAYOUT     = "2006-01-02 15:04:05.999"
	DEFAUL_SKIP             = 3
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type Formatter interface {
	Format(b *bytes.Buffer, level string, format string, args ...interface{}) []byte
	FormatWithDetail(b *bytes.Buffer, level string, format string, args ...interface{}) []byte
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

func (this *DefaultFormatter) FormatWithDetail(b *bytes.Buffer, level string, format string, args ...interface{}) []byte {
	msg := strings.Replace(this.Template, ASCTIME_MARK, this.FormatTime(), -1)
	msg = strings.Replace(msg, LEVELNAME_MARK, level, -1)
	msg = strings.Replace(msg, MSG_MARK, fmt.Sprintf(format, args...), -1)

	b.WriteString(msg)
	b.WriteByte('\n')
	b.Write(stack(DEFAUL_SKIP))
	b.WriteByte('\n')
	return b.Bytes()
}

func (this *DefaultFormatter) FormatTime() string {
	return time.Now().Format(this.TimeLayout)
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
