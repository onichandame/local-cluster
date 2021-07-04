package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/job"
	"time"
)

func main() {
	preBootstrap()

	r := gin.Default()
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			"ok":        true,
		})
	})
	r.Run()
}

func preBootstrap() {
	config.ConfigInit()
	db.DBInit()
	job.JobInit()
}
