package main

import (
	"github.com/baa-middleware/logger"
	"github.com/baa-middleware/recovery"
	"github.com/go-baa/baa"
)

func hello(c *baa.Context) {
	c.String(200, "Hello, World!\n")
}

func main() {
	b := baa.New()
	b.Use(logger.Logger())
	b.Use(recovery.Recovery())

	b.Get("/", hello)

	b.Run(":8001")
}
