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

var roomsCode = make(map[string]string)
var codeLock = sync.Mutex{}

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

		go UpdateCodeLock(socketReq.RoomID, socketReq.SourceCode)

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
	go CodeSaver()

	r.Run()
}

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func SaveSpecificRoom(roomID string) {
	codeLock.Lock()
	defer codeLock.Unlock()

	if _, ok := roomsCode[roomID]; ok {
		room.SaveRoomObject(room.RoomObject{
			ID:         roomID,
			SourceCode: roomsCode[roomID],
			Language:   "",
		})
	}
}

func UpdateCodeLock(roomID, sourceCode string) {
	codeLock.Lock()
	defer codeLock.Unlock()

	roomsCode[roomID] = sourceCode
}

func CodeSaver() {
	for {
		func() {
			time.Sleep(time.Second * 10)
			codeLock.Lock()
			defer codeLock.Unlock()

			for roomID, sourceCode := range roomsCode {
				fmt.Println("saved code for room: ", roomID)
				room.SaveRoomObject(room.RoomObject{
					ID:         roomID,
					SourceCode: sourceCode,
					Language:   "",
				})
			}

			roomsCode = make(map[string]string)
		}()
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
