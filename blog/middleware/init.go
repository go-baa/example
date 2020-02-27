package middleware

import (
	"github.com/baa-middleware/gzip"
	"github.com/baa-middleware/recovery"
	"github.com/baa-middleware/session"
	"github.com/go-baa/baa"
)

// Initializes 初始化中间件
func Initializes(b *baa.Baa) {
	// Recovery
	b.Use(recovery.Recovery())

	// Session
	cacheOptions := session.MemoryOptions{}
	cacheOptions.BytesLimit = 32 * 1024 * 1024 // 32M

	b.Use(session.Middleware(session.Options{
		Name: "BaaBlogSession",
		Provider: &session.ProviderOptions{
			Adapter: "memory",
			Config:  cacheOptions,
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
