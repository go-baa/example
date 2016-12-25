package accesslog

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/baa.v1"
)

const (
	// FormatterSymbolStart 日志格式开始符号
	FormatterSymbolStart = '%'
	// FormatterSymbolEnd 日志格式结束符号
	FormatterSymbolEnd = '%'
	// DefaultFormatter 默认日志格式
	DefaultFormatter = `%remote_addr% "%http_x_forwarded_for%" %request% %status% %body_bytes_sent% %exec_time% %http_referer% %http_user_agent%`
)

// formatter 日志格式处理器
type formatter struct {
	Text string
	Vars []string
}

// newFormatter 创建一个日志格式示例
func newFormatter(format string) *formatter {
	f := new(formatter)
	f.perpare(format)
	return f
}

// perpare 预处理日志格式字符串
func (f *formatter) perpare(format string) {
	fBuf := new(bytes.Buffer)
	vBuf := new(bytes.Buffer)
	l := len(format)
	for i := 0; i < l; i++ {
		if c := format[i]; c != FormatterSymbolStart {
			fBuf.WriteByte(c)
			continue
		}

		vBuf.Reset()
		i += len(string(FormatterSymbolStart))
		for {
			if format[i] == FormatterSymbolEnd {
				break
			}
			vBuf.WriteByte(format[i])
			i++
		}
		if vBuf.Len() > 0 {
			f.Vars = append(f.Vars, vBuf.String())
			fBuf.Write([]byte("%v"))
		}
	}
	fBuf.WriteString("\n")
	f.Text = fBuf.String()
}

// build 构建日志行
func (f *formatter) build(c *baa.Context, start time.Time) string {
	var buf []interface{}
	for _, v := range f.Vars {
		switch v {
		case "hostname":
			hostname, _ := os.Hostname()
			buf = append(buf, hostname)
		case "time_iso8601":
			buf = append(buf, start.Format("2006-01-02T15:04:05+0800"))
		case "query_string":
			var qBuf bytes.Buffer
			for key, value := range c.Querys() {
				qBuf.WriteString(fmt.Sprintf("%s=%v&", key, value))
			}
			buf = append(buf, strings.TrimRight(qBuf.String(), "&"))
		case "http_host":
			buf = append(buf, c.Req.Host)
		case "exec_time":
			buf = append(buf, time.Since(start).String())
		case "method":
			buf = append(buf, c.Req.Method)
		case "remote_addr":
			buf = append(buf, c.RemoteAddr())
		case "request":
			buf = append(buf, c.Req.Method+" "+c.Req.RequestURI+" "+c.Req.Proto)
		case "request_uri":
			buf = append(buf, c.Req.RequestURI)
		case "status":
			buf = append(buf, c.Resp.Status())
		case "status_text":
			buf = append(buf, http.StatusText(c.Resp.Status()))
		case "body_bytes_sent":
			buf = append(buf, c.Resp.Size())
		default:
			if strings.HasPrefix(v, "http_") {
				buf = append(buf, c.Req.Header.Get(strings.Replace(string(v[5:]), "_", "-", -1)))
			} else {
				buf = append(buf, "")
			}
		}
	}
	return fmt.Sprintf(f.Text, buf...)
}
