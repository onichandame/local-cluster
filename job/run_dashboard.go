package job

import (
	"github.com/onichandame/local-cluster/application"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/instance_group"
)

var runDashboard = job{
	immediate: true,
	name:      "RunDashboard",
	dependsOn: []*job{&auditInstances},
	run: func() error {
		ig, err := getOrCreateInsGrp()
		if err != nil {
			return err
		}
		if err := instancegroup.Start(ig); err != nil {
			return err
		}
		return nil
	},
}

func getOrCreateApp() (*model.Application, error) {
	var err error
	app := model.Application{}
	app.Name = "dashboard"
	app.Interfaces = []model.ApplicationInterface{
		{
			Name:      "main",
			PortByEnv: "PORT",
		},
	}
	app.Specs = []model.ApplicationSpec{
		{
			DownloadUrl: "https://github.com/onichandame/local-cluster-dashboard/releases/download/latest/release.tar.gz",
			Platform:    "linux",
			Arch:        "amd64",
			Command:     "npx",
			Args:        "serve build"},
	}
	if err := application.Prepare(&app); err != nil {
		return nil, err
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
