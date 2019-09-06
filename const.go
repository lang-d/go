package golog

// 日志级别相关
const (
	DEBUG = iota
	INFO
	WARN
	FATAL
)

// 日志文件分割相关
const (
	SPLIT_BY_HOUR      = "split_by_hour"
	SPLIT_BY_DAY       = "split_by_day"
	SPLIT_BY_SIZE      = "split_by_size"
	DEFAULT_BACKUPN    = 10
	DEFAULT_SPLIT_SIZE = 100 * 1024 * 1024 // 100M

)
