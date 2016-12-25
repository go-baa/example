// Package accesslog 提供了一个baa的访问日志中间件
package accesslog

import (
	"fmt"
	"time"

	"gopkg.in/baa.v1"
)

// Logger 访问日志接口
type Logger interface {
	Log(string)
	Flush() error
	Config(*Options) error
}

// Options 访问日志配置
type Options struct {
	Open    bool
	Adapter string
	Format  string
	Config  map[string]interface{}
}

// Get 从配置中获取一个值
func (o Options) Get(key string, defaultValue interface{}) interface{} {
	v, exists := o.Config[key]
	if !exists {
		return defaultValue
	}
	return v
}

type instanceFunc func() Logger

var adapters = make(map[string]instanceFunc)

// New 创建一个新的访问记录中间件
func New(o Options) baa.Middleware {
	if !o.Open {
		return func(c *baa.Context) { c.Next() }
	}

	if len(o.Adapter) == 0 {
		panic("accesslog.New: cannot use empty adapter")
	}

	// adapter
	f, ok := adapters[o.Adapter]
	if !ok {
		panic("accesslog.New: unknown adapter (forgot to import?)")
	}
	adapter := f()

	// Set default logger format
	if o.Format == "" {
		o.Format = DefaultFormatter
	}

	if err := adapter.Config(&o); err != nil {
		panic(fmt.Sprintf("accesslog.New: %s incorrect configuration, %s", o.Adapter, err.Error()))
	}

	// new formater
	formater := newFormatter(o.Format)

	return func(c *baa.Context) {
		start := time.Now()

		c.Next()

		// 异步日志
		go func() {
			line := formater.build(c, start)
			adapter.Log(line)
		}()
	}
}

// Register 注册新的适配器
func Register(name string, adapter instanceFunc) {
	if adapter == nil {
		panic("accesslog.Register: cannot register adapter with nil func")
	}
	if _, ok := adapters[name]; ok {
		panic(fmt.Errorf("accesslog.Register: cannot register adapter '%s' twice", name))
	}
	adapters[name] = adapter
}
