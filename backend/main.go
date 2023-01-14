package main

import (
	"codebooks/judge"
	"codebooks/room"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	apiGroup := r.Group("/api")
	apiGroup.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
	})

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

	r.Run()
}
