package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/baa-middleware/logger"
	"github.com/baa-middleware/recovery"
	"github.com/go-baa/baa"
)

func main() {
	b := baa.New()
	b.Use(logger.Logger())
	b.Use(recovery.Recovery())

	b.Get("/", func(c *baa.Context) {
		c.String(200, "Hello, World!\n")
	})

	b.Post("/", func(c *baa.Context) {
		defer c.Req.Body.Close()
		body, err := ioutil.ReadAll(c.Req.Body)
		if err != nil {
			c.Error(err)
			return
		}
		var data map[string]interface{}
		body = []byte(strings.Replace(string(body), "\\x22", "\"", -1))
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.Error(err)
			return
		}
		for k, v := range data {
			c.Resp.Write([]byte(fmt.Sprintf("%s: %v\n", k, v)))
		}
	})

	b.Get("/tpl", func(c *baa.Context) {
		c.Set("name", "micate")
		html, _ := c.Fetch("template/test.html")
		c.Text(200, html)
	})

	b.Get("/file", func(c *baa.Context) {
		c.HTML(200, "template/upload.html")
	})

	b.Post("/file", func(c *baa.Context) {
		c.Set("posts", c.Posts())
		c.Set("a", c.Query("a"))
		_, fh, _ := c.GetFile("file1")
		c.Set("file", fh.Filename)
		c.JSON(200, c.Gets())
	})

	b.Run(":8001")
}
