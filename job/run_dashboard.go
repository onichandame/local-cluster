package job

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/application"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/instance_group"
)

const (
	interfaceName = "main"
)

var runDashboard = job{
	immediate: true,
	name:      "RunDashboard",
	dependsOn: []*job{&auditInstances, &initInterfaces, &initProxyManager, &initCacheManager, &auditEntrances},
	run: func() (err error) {
		err = nil
		defer func() {
			if er := recover(); er != nil {
				if e, ok := er.(error); ok {
					err = e
				}
			}
		}()
		if ig, err := getOrCreateInsGrp(); err != nil {
			panic(err)
		} else {
			if err := instancegroup.Start(ig); err != nil {
				panic(err)
			}
			if err := createEntrance(ig); err != nil {
				panic(err)
			}
		}
		return err
	},
}

func getOrCreateApp() (*model.Application, error) {
	appName := "dashboard"
	var app model.Application
	var err error
	if err = db.Db.Where("name = ?", appName).First(&app).Error; err == nil {
		return &app, err
	}
	app.Name = "dashboard"
	app.Interfaces = []model.ApplicationInterface{
		{
			Name:      interfaceName,
			PortByEnv: "PORT",
		},
	}
	app.Specs = []model.ApplicationSpec{
		{
			DownloadUrl: "https://github.com/onichandame/local-cluster-dashboard/releases/download/latest/release.tar.gz",
			Platform:    "linux",
			Arch:        "amd64",
			Command:     "npx",
			Args:        "serve -s build"},
	}
	if err = application.Prepare(&app); err != nil {
		return &app, err
	}
	return &app, err
}

func getOrCreateInsGrp() (*model.InstanceGroup, error) {
	var err error
	app, err := getOrCreateApp()
	if err != nil {
		return nil, err
	}
	ig := model.InstanceGroup{}
	ig.Replicas = 2
	ig.ApplicationID = app.ID
	if err := db.Db.Where("application_id = ?", ig.ApplicationID).FirstOrCreate(&ig).Error; err != nil {
		return nil, err
	}
	return &ig, err
}

func createEntrance(igDef *model.InstanceGroup) (err error) {
	err = nil
	defer func() {
		if e := recover(); e != nil {
			if er, ok := e.(error); ok {
				err = er
			}
		}
	}()
	ent := new(model.Entrance)
	if err = db.Db.Preload("Interfaces.Definition").First(igDef, igDef.ID).Error; err != nil {
		panic(err)
	}
	var igIf model.ServiceInterface
	for _, i := range igDef.Interfaces {
		if i.Definition.Name == interfaceName {
			igIf = i
		}
	}
	if igIf.ID == 0 {
		panic(errors.New(fmt.Sprintf("failed to find interface for dashboard! cannot create entrance")))
	}
	ent.Backend = igIf
	ent.Name = "dashboard"
	if err = db.Db.Create(ent).Error; err != nil {
		panic(err)
	}
	return err
}
