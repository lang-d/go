package golog

var (
	log = NewLogger()
)

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func AddHandler(handler Handler) {
	log.AddHandler(handler)
}

func SetHandlers(handlers []Handler) {
	log.SetHandlers(handlers)
}

func SetLevel(lev *level) {
	log.SetLevel(lev)
}

func SetHandler(handler Handler) {
	log.SetHandlers([]Handler{handler})
}
