package middleware

import (
	"github.com/go-baa/example/blog/modules/setting"

	"github.com/baa-middleware/gzip"
	"github.com/baa-middleware/recovery"
	"github.com/baa-middleware/session"

	"gopkg.in/baa.v1"
)

// Initializes 初始化中间件
func Initializes(b *baa.Baa) {
	// Recovery
	b.Use(recovery.Recovery())

	// Session
	redisOptions := session.RedisOptions{}
	redisOptions.Addr = setting.Config.MustString("session.redis.addr", "")
	redisOptions.DB = setting.Config.MustInt64("session.redis.db", 0)
	redisOptions.Password = setting.Config.MustString("session.redis.password", "")
	redisOptions.Prefix = setting.Config.MustString("session.redis.prefix", "")

	b.Use(session.Middleware(session.Options{
		Name: "BaaBlogSession",
		Provider: &session.ProviderOptions{
			Adapter: "redis",
			Config:  redisOptions,
		},
		Cookie: &session.CookieOptions{
			Path:     "/",
			HttpOnly: true,
			LifeTime: 86400,
		},
		MaxLifeTime: 86400,
	}))

	// Gzip
	if baa.Env == baa.PROD {
		b.Use(gzip.Gzip(gzip.Options{CompressionLevel: 4}))
	}
}
