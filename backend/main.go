package main

import (
	"codebooks/judge"
	"codebooks/room"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func main() {
	r := gin.Default()

	// socket server
	server := socketio.NewServer(nil)
	server.OnConnect("/", func(c socketio.Conn) error {
		c.SetContext("")
		fmt.Println("connected:", c.ID())

		return nil
	})

	// API server
	r.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	apiGroup := r.Group("/api")

	// v1 API
	apiV1 := apiGroup.Group("/v1")
	apiV1.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"api":    "v1",
		})
	})

	// rooms
	roomAPI := apiV1.Group("/room")
	roomAPI.POST("/create", room.CreateRoom)
	roomAPI.GET("/get", room.GetRoom)

	// judge
	judgeAPI := apiV1.Group("/judge")
	judgeAPI.POST("/submit", judge.Submit)

	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))
	go server.Serve()
	defer server.Close()

	r.Run()
}
