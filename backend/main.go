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
	r.Use(GinMiddleware("http://localhost:3000"))

	// socket server
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(c socketio.Conn) error {
		c.SetContext("")
		fmt.Println("connected:", c.ID())
		return nil
	})

	server.OnEvent("/", "result", func(c socketio.Conn, msg string) {
		fmt.Println(fmt.Printf("received msg %s from id %s", msg, c.ID()))
	})

	server.OnEvent("/", "joinroom", func(c socketio.Conn, msg string) {
		fmt.Println("joining room", msg)
		c.Join(msg)
	})

	// API server
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
	roomAPI.POST("/save", room.SaveRoom)
	roomAPI.GET("/check", room.CheckRoom)

	// judge
	judgeAPI := apiV1.Group("/judge")
	judgeAPI.POST("/submit", judge.Submit)

	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))

	go server.Serve()
	defer server.Close()

	r.Run()
}

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}
