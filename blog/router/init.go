package router

import (
	"github.com/go-baa/baa"
	"github.com/go-baa/cache"
	_ "github.com/go-baa/cache/redis"
	"github.com/go-baa/example/blog/modules/template"
	"github.com/go-baa/log"
	"github.com/go-baa/render"
	"github.com/go-baa/setting"
)

// Initializes 初始化
func Initializes(b *baa.Baa) {
	// 日志
	b.SetDI("logger", log.Default())

	// 静态目录
	b.Static("/assets/", "public/assets/", false, nil)
	b.Static("/upload/", "public/upload/", false, nil)
	b.Static("/favicon.ico", "public/favicon.ico", false, nil)

	// 模板渲染
	b.SetDI("render", render.New(render.Options{
		Baa:        b,
		Root:       "template/",
		Extensions: []string{".html", ".tmpl"},
		FuncMap:    template.Funcs(b),
	}))

	// cache
	b.SetDI("cache", cache.New(cache.Options{
		Name:    setting.Config.MustString("cache.section", ""),
		Prefix:  setting.Config.MustString("cache.prefix", ""),
		Adapter: setting.Config.MustString("cache.adapter", "memory"),
		Config: map[string]interface{}{
			"host":       setting.Config.MustString("cache.redis.host", "127.0.0.1"),
			"port":       setting.Config.MustString("cache.redis.port", "6379"),
			"password":   setting.Config.MustString("cache.redis.password", ""),
			"poolsize":   setting.Config.MustInt("cache.redis.poolsize", 10),
			"bytesLimit": int64(32 * 1024 * 1024),
		},
	}))
}
