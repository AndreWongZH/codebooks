package main

import (
	"codebooks/judge"
	"codebooks/room"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type SocketStruct struct {
	SourceCode string `json:"source_code"`
	User       string `json:"user"`
	RoomID     string `json:"room_id"`
}

var roomsList = make(map[string][]string)
var roomsLock = sync.Mutex{}

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

	server.OnEvent("/", "pongpong", func(c socketio.Conn, msg string) {
		fmt.Printf("receive pong\n")
		var socketReq SocketStruct
		json.Unmarshal([]byte(msg), &socketReq)
		go AddUserToRoom(socketReq.User, socketReq.RoomID)
		fmt.Printf("receive pong from %v  of room %v\n", socketReq.User, socketReq.RoomID)
	})

	server.OnEvent("/", "joinroom", func(c socketio.Conn, msg string) {
		var socketReq SocketStruct
		json.Unmarshal([]byte(msg), &socketReq)
		fmt.Printf("%v joining room %v\n", socketReq.User, socketReq.RoomID)
		go AddUserToRoom(socketReq.User, socketReq.RoomID)
		c.Join(socketReq.RoomID)
	})

	server.OnEvent("/", "edit", func(c socketio.Conn, msg string) {
		var socketReq SocketStruct
		json.Unmarshal([]byte(msg), &socketReq)
		fmt.Println(socketReq)

		b, _ := json.Marshal(&SocketStruct{
			SourceCode: socketReq.SourceCode,
			User:       socketReq.User,
			RoomID:     socketReq.RoomID,
		})
		server.BroadcastToRoom("/", socketReq.RoomID, "newcode", string(b))
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

	go ActiveRoomPinger(server)

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

func AddUserToRoom(userID, roomID string) {
	roomsLock.Lock()
	defer roomsLock.Unlock()
	for _, userIn := range roomsList[roomID] {
		if userIn == userID {
			return
		}
	}
	roomsList[roomID] = append(roomsList[roomID], userID)
}

func ActiveRoomPinger(server *socketio.Server) {
	for {
		func() {
			roomsLock.Lock()
			for k := range roomsList {
				roomsList[k] = roomsList[k][0:0]
				server.BroadcastToRoom("/", k, "ping")
			}
			defer roomsLock.Unlock()
		}()
		time.Sleep(time.Second * 10)

		roomsLock.Lock()
		for k, v := range roomsList {
			server.BroadcastToRoom("/", k, "active_users", v)
		}
		roomsLock.Unlock()
		time.Sleep(time.Second * 5)
	}
}
