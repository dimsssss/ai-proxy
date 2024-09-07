package main

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dimsssss/ai-proxy/cmd/database"
	"github.com/dimsssss/ai-proxy/cmd/env"
	"github.com/dimsssss/ai-proxy/internal/ratelimit"
)

func main() {

	env.LoadEnv()

	db := database.Connection()
	ratelimits := ratelimit.NewRateLimitWith(ratelimit.GetServices(db))
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.GET("/api/v1/llm", func(c *gin.Context) {
		// TODO require bearer jwt token that include service id
		service, ok := ratelimits.Services[c.Query("id")]

		if !ok {
			c.JSON(400, gin.H{
				"message": "not exist service",
			})
			return
		}

		success, err := ratelimits.ProcessRateLimit(time.Now(), service)

		if success {
			bodysize := ratelimit.GetLlmResult()
			service.Update(bodysize)

			c.JSON(200, gin.H{
				"message": "success",
			})
		}

		if err != nil {
			switch e := err.(type) {
			case *ratelimit.RateLimitError:
				c.JSON(429, gin.H{
					"message": e.Error(),
				})
				return
			case *ratelimit.InvalidError:
				c.JSON(400, gin.H{
					"message": err.Error(),
				})
				return
			}
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
