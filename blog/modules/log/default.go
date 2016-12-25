package log

import (
	"log"
	"os"

	"git.vodjk.com/golang/common/modules/log/report"
	"git.vodjk.com/golang/common/modules/setting"
)

// _defaultLogger 日志接口
var _defaultLogger Logger

// Default 返回默认日志器
func Default() Logger {
	return _defaultLogger
}

func Level() LogLevel {
	return _defaultLogger.Level()
}

func SetLevel(l LogLevel) {
	_defaultLogger.SetLevel(l)
}

func Flush() {
	_defaultLogger.Flush()
}

func Print(v ...interface{}) {
	_defaultLogger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	_defaultLogger.Printf(format, v...)
}

func Println(v ...interface{}) {
	_defaultLogger.Println(v...)
}

func Fatal(v ...interface{}) {
	_defaultLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	_defaultLogger.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	_defaultLogger.Fatalln(v...)
}

func Panic(v ...interface{}) {
	_defaultLogger.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	_defaultLogger.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	_defaultLogger.Panicln(v...)
}

func Error(v ...interface{}) {
	_defaultLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	_defaultLogger.Errorf(format, v...)
}

func Errorln(v ...interface{}) {
	_defaultLogger.Errorln(v...)
}

func Warn(v ...interface{}) {
	_defaultLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	_defaultLogger.Warnf(format, v...)
}

func Warnln(v ...interface{}) {
	_defaultLogger.Warnln(v...)
}

func Info(v ...interface{}) {
	_defaultLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	_defaultLogger.Infof(format, v...)
}

func Infoln(v ...interface{}) {
	_defaultLogger.Infoln(v...)
}

func Debug(v ...interface{}) {
	_defaultLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	_defaultLogger.Debugf(format, v...)
}

func Debugln(v ...interface{}) {
	_defaultLogger.Debugln(v...)
}

func init() {
	var f *os.File
	file := setting.Config.MustString("log.file", "")
	if file == "" || file == "os.Stderr" {
		f = os.Stderr
	} else if file == "os.Stdout" {
		f = os.Stdout
	} else {
		var err error
		f, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening logfile: %v\n", err)
		}
		// 不能关闭文件
		//defer f.Close()
	}
	_defaultLogger = New(f, "["+setting.AppName+"] ", log.LstdFlags)
	_defaultLogger.SetLevel(LogLevel(setting.Config.MustInt("log.level", int(LOG_WARN))))
	_defaultLogger.(*tLogger).SetCallerLevel(3)

	if setting.Config.MustBool("log.report.open", false) {
		// report Error, Exception
		flume, err := report.NewFlume(setting.Config.MustString("log.report.host", "127.0.0.1"), setting.Config.MustString("log.report.port", "6000"))
		if err == nil {
			_defaultLogger.(*tLogger).report = newReport(flume)
		}
	}
}
