package main

import (
	"github.com/baa-middleware/accesslog"
	"github.com/baa-middleware/recovery"
	"github.com/go-baa/baa"
)

func hello(c *baa.Context) {
	c.String(200, "Hello, World!\n")
}

func main() {
	b := baa.New()
	b.Use(accesslog.Logger())
	b.Use(recovery.Recovery())

	b.Get("/", hello)
	b.Get("/error", func(c *baa.Context) {
		panic("this is panic test")
	})
	b.Get("/debug", func(c *baa.Context) {
		debug := c.QueryBool("debug")
		c.Baa().SetDebug(debug)
		if debug {
			c.String(200, "debug open")
		} else {
			c.String(200, "debug off")
		}
	})

	b.Run(":1323")
}
