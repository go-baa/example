package main

import (
	"github.com/baa-middleware/accesslog"
	"github.com/baa-middleware/recovery"
	"gopkg.in/baa.v1"
)

func hello(c *baa.Context) {
	c.String(200, "Hello, World!\n")
}

func main() {
	b := baa.New()
	b.Use(accesslog.Logger())
	b.Use(recovery.Recovery())

	b.Get("/", hello)

	b.Run(":1323")
}
