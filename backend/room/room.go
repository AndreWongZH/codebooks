package room

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRoom(c *gin.Context) {
	// generate random base64 room_id
	b := make([]byte, 18)
	rand.Read(b)
	room_id := base64.URLEncoding.EncodeToString(b)

	// add to storage object (TODO)

	// response
	c.JSON(http.StatusOK, gin.H{
		"room_id": room_id,
	})
}

func GetRoom(c *gin.Context) {
	// response
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
