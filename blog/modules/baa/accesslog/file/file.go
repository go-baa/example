package file

import (
	"os"
	"strings"
	"sync"

	"github.com/go-baa/example/blog/modules/baa/accesslog"
)

const (
	// DefaultSize 日志缓冲行数
	DefaultSize = 10000
)

// File a adapter of accesslog
type File struct {
	m    *sync.RWMutex
	f    *os.File
	size int
	buf  []string
}

// New create a accesslog instance with file
func New() accesslog.Logger {
	return new(File)
}

// Log 记录一条日志
func (l *File) Log(line string) {
	// append && flush
	go func() {
		l.m.Lock()
		l.buf = append(l.buf, line)
		n := len(l.buf)
		l.m.Unlock()
		if n > l.size {
			l.Flush()
		}
	}()
}

// Flush 立即写入缓存信息
func (l *File) Flush() error {
	l.m.Lock()
	defer l.m.Unlock()
	l.f.WriteString(strings.Join(l.buf, ""))
	if err := l.f.Sync(); err != nil {
		return err
	}
	l.buf = l.buf[:0]
	return nil
}

// Config 初始化日志
func (l *File) Config(o *accesslog.Options) error {
	file := o.Get("file", "").(string)
	if file == "" || file == "os.Stderr" {
		l.f = os.Stderr
	} else if file == "os.Stdout" {
		l.f = os.Stdout
	} else {
		var err error
		l.f, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			return err
		}
	}

	l.size = o.Get("size", DefaultSize).(int)
	if l.size == 0 {
		l.size = DefaultSize
	}
	l.m = &sync.RWMutex{}
	l.buf = make([]string, l.size)

	return nil
}

func init() {
	accesslog.Register("file", New)
}
