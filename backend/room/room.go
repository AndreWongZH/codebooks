package room

import (
	// "crypto/rand"

	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CheckRoom(c *gin.Context) {
	roomID := c.Query("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "specify room_id",
		})
		return
	}

	exists := CheckRoomExistsObject(roomID)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"result": exists,
	})
}

var roomids []int

func CreateRoom(c *gin.Context) {
	// generate random base64 room_id
	// b := make([]byte, 18)
	// rand.Read(b)
	// room_id := base64.URLEncoding.EncodeToString(b)
	for {
		room_id := rand.Intn(8999) + 1000

		var result bool = false
		for _, x := range roomids {
			if x == room_id {
				result = true
				break
			}
		}

		if !result {
			roomids = append(roomids, room_id)

			// add to storage object (TODO)
			room, err := ReadRoomObject(strconv.Itoa(room_id))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
				})
				return
			}

			// response
			c.JSON(http.StatusOK, gin.H{
				"room_id":     room.ID,
				"source_code": room.SourceCode,
				"language":    room.Language,
			})

			break
		}
	}

}

func GetRoom(c *gin.Context) {
	// read the room file from storage object
	roomID := c.Query("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "no room id",
		})
		return
	}

	room, err := ReadRoomObject(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "file read error",
		})
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"result": room,
	})
}

func SaveRoom(c *gin.Context) {
	var req SaveRoomRequest

	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "error parsing request body",
		})
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error unmarshalling json",
		})
	}

	SaveRoomObject(RoomObject{
		ID:         req.RoomID,
		SourceCode: req.SourceCode,
		Language:   req.Language,
	})

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

type SaveRoomRequest struct {
	RoomID     string `json:"room_id"`
	SourceCode string `json:"source_code"`
	Language   string `json:"language"`
}
