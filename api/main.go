package main

import (
	"github.com/baa-middleware/accesslog"
	"github.com/baa-middleware/gzip"
	"github.com/baa-middleware/recovery"
	"github.com/go-baa/example/api/controller"
	"gopkg.in/baa.v1"
)

func main() {
	app := baa.Default()
	app.Use(recovery.Recovery())
	app.Use(accesslog.Logger())

	// production mode enable gzip
	if !app.Debug() {
		app.Use(gzip.Gzip(gzip.Options{
			CompressionLevel: 9,
		}))
	}

	app.Get("/", controller.IndexController.Index)
	app.Get("/list/:page", controller.IndexController.Index)
	app.Group("/show", func() {
		app.Get("/", func(c *baa.Context) {
			c.Redirect(302, "/")
		})
		app.Get("/:id", controller.IndexController.Show)
	})

	app.Run(":1323")
}
