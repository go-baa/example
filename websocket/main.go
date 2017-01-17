package main

import (
	"fmt"
	"time"

	"github.com/go-baa/baa"
	"github.com/gorilla/websocket"
)

func accesslog() baa.Middleware {
	return func(c *baa.Context) {
		start := time.Now()

		c.Next()

		fmt.Printf("[websocket] %s %s %v %v\n", c.Req.Method, c.Req.RequestURI, c.Resp.Status(), time.Since(start))
	}
}

func main() {
	fmt.Println("begin")
	app := baa.Default()
	app.Use(accesslog())

	app.Static("/", ".", true, nil)
	app.Websocket("/socket", func(ws *websocket.Conn) {
		ws.SetPingHandler(func(appData string) error {
			fmt.Println("websocket ping: ", appData)
			return nil
		})
		ws.SetPongHandler(func(appData string) error {
			fmt.Println("websocket pong: ", appData)
			return nil
		})
		ws.SetCloseHandler(func(code int, text string) error {
			fmt.Println("websocket close: ", code, text)
			return nil
		})
		for {
			fmt.Println("websocket retry read...")
			go func() {
				time.Sleep(time.Second * 10)
				ws.Close()
			}()
			messageType, data, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err) {
					app.Logger().Println("websocket ReadMessage error: connection is closed")
				} else {
					app.Logger().Println("websocket ReadMessage error:", err)
				}
				ws.Close()
				return
			}
			fmt.Println("websocket receive: ", messageType, string(data))
			err = ws.WriteMessage(messageType, data)
			if err != nil {
				app.Logger().Println("websocket WriteMessage error:", err)
				return
			}
		}
	})

	app.Run(":1234")

	fmt.Println("end")
}
