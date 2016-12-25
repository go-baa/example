package main

import (
	"github.com/go-baa/example/blog/middleware"
	"github.com/go-baa/example/blog/modules/setting"
	"github.com/go-baa/example/blog/router"

	"gopkg.in/baa.v1"
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
	app.Run(setting.Config.MustString("http.address", "0.0.0.0") + ":" + setting.Config.MustString("http.port", "8001"))
}
