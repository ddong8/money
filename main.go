package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

func main() {
	r := gin.Default()

	mon := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitoring/tasks", // RootPath specifies the root for asynqmon app
		RedisConnOpt: asynq.RedisClientOpt{Addr: "redis:6379"},
	})
	r.Any("/monitoring/tasks/*a", gin.WrapH(mon))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
