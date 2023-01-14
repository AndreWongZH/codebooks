package main

import (
	"codebooks/judge"
	"codebooks/room"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type LastEdit struct {
	Timestamp    time.Time
	LastEditUser string
}

type SocketStruct struct {
	SourceCode string `json:"source_code"`
	User       string `json:"user"`
	RoomID     string `json:"room_id"`
}

var rooms = map[string]LastEdit{}

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

	server.OnEvent("/", "edit", func(c socketio.Conn, msg string) {
		var socketReq SocketStruct
		json.Unmarshal([]byte(msg), &socketReq)
		fmt.Println(socketReq)
		c.Join(socketReq.RoomID)
		fmt.Println(c.Rooms())

		b, _ := json.Marshal(&SocketStruct{
			SourceCode: socketReq.SourceCode,
			User:       socketReq.User,
			RoomID:     socketReq.RoomID,
		})
		server.BroadcastToRoom("/", socketReq.RoomID, "newcode", string(b))

		// emit back to user that he cannot edit
		c.Emit("edit")

		// roomID, user := "", ""
		// timestamp := time.Now()

		// if _, ok := rooms[roomID]; !ok {
		// 	rooms[roomID] = LastEdit{
		// 		Timestamp:    time.Unix(0, 0),
		// 		LastEditUser: "",
		// 	}
		// }

		// if rooms[roomID].Timestamp.Unix()-timestamp.Unix() < 5 && rooms[roomID].LastEditUser != user {
		// 	// error if user is different
		// 	c.Emit("noedit")
		// } else {
		// 	// allow the edit
		// 	rooms[roomID] = LastEdit{
		// 		LastEditUser: user,
		// 		Timestamp:    timestamp,
		// 	}
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
