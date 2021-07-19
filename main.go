package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/application"
	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/instance"
	"github.com/onichandame/local-cluster/job"
	"github.com/onichandame/local-cluster/route"
)

func main() {
	preBootstrap()

	r := gin.Default()
	route.RoutesInit(r)
	r.Run()
}

func preBootstrap() {
	config.Init()
	db.Init()
	job.JobsInit()
	if err := application.AuditCache(); err != nil {
		panic(err)
	}

	if err := instance.Audit(); err != nil {
		panic(err)
	}
}
