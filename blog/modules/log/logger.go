package log

// LogLevel 日志级别
type LogLevel int

const (
	// !nashtsai! following level also match syslog.Priority value
	LOG_UNKNOWN LogLevel = iota - 2 // -2
	LOG_OFF     LogLevel = iota - 1 // 0
	LOG_FATAL                       // 1
	LOG_PANIC                       // 2
	LOG_ERROR   LogLevel = iota + 1 // 5
	LOG_WARN                        // 6
	LOG_INFO    LogLevel = iota + 4 // 10
	LOG_DEBUG                       // 11
)

var levelNames = map[LogLevel]string{
	LOG_UNKNOWN: "UNKNOWN",
	LOG_FATAL:   "FATAL",
	LOG_PANIC:   "PANIC",
	LOG_ERROR:   "ERROR",
	LOG_WARN:    "WARN",
	LOG_INFO:    "INFO",
	LOG_DEBUG:   "DEBUG",
}

// Logger interface
type Logger interface {
	Level() LogLevel
	SetLevel(l LogLevel)
	Flush()
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
}
