package main

import (
	"github.com/go-baa/baa"
	"github.com/go-baa/cache"
)

func main() {
	// new app
	app := baa.New()

	// register cache
	app.SetDI("cache", cache.New(cache.Options{
		Name:     "cache",
		Adapter:  "memory",
		Config:   map[string]string{},
		Interval: 60,
	}))

	// router
	app.Get("/", func(c *baa.Context) {
		ca := c.DI("cache").(cache.Cacher)
		ca.Set("test", "baa", 10)
		v := ca.Get("test").(string)
		c.String(200, v)
	})

	// run app
	app.Run(":1323")
}
