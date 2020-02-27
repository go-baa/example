package main

import (
	"github.com/go-baa/baa"
	"github.com/go-baa/example/blog/middleware"
	"github.com/go-baa/example/blog/router"
	"github.com/go-baa/setting"
)

// main 入口
func main() {
	app := baa.Default()

	// 初始化middleware
	middleware.Initializes(app)

	// 初始化路由
	router.Initializes(app)
	router.Router(app)

	// 运行
	addr := setting.Config.MustString("http.address", "0.0.0.0")
	port := setting.Config.MustString("http.port", "1323")
	app.Run(addr + ":" + port)
}
