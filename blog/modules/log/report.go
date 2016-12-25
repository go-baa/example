package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

// mReport 错误报告结构
type mReport struct {
	out io.Writer
}

// 创建一个Report结构
// out 是一个io.Writer的实现
func newReport(out io.Writer) *mReport {
	return &mReport{out: out}
}

func (r *mReport) write(appName string, appVersion string, file string, line int, level LogLevel, msg string) {
	hostname, _ := os.Hostname()
	r.output(fmt.Sprintf(`{"file":"%s", "line":"%d", "level": "%v", "levelname": "%s", "message": "%s", "timestamp": "%s", "hostname": "%s", "appname":"%s", "appversion": "%s"}`, file, line, level, _defaultLogger.(*tLogger).LevelName(level), msg, time.Now().Format("2006-01-02T15:04:05+0800"), hostname, appName, appVersion))
}

func (r *mReport) output(line string) {
	r.out.Write([]byte(line))
}
