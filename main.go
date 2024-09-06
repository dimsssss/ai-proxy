package main

import (
	"github.com/gin-gonic/gin"

	"github.com/dimsssss/ai-proxy/internal/database"
	"github.com/dimsssss/ai-proxy/internal/env"
)

func main() {

	env.LoadEnv()

	database.Connection()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
