package job

import (
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/instance"
)

var runDashboard = job{
	immediate: true,
	name:      "RunDashboard",
	dependsOn: []*job{&auditInstances},
	run: func() error {
		ins, err := getOrCreateIns()
		if err != nil {
			return err
		}
		if err := instance.RunInstance(ins); err != nil {
			return err
		}
		return nil
	},
}

func getOrCreateApp() (*model.Application, error) {
	var err error
	app := model.Application{}
	app.Name = "dashboard"
	app.Specs = []model.ApplicationSpec{{DownloadUrl: "https://github.com/onichandame/local-cluster-dashboard/releases/download/latest/release.tar.gz", Platform: "linux", Arch: "amd64", Entrypoint: "npx", Args: "serve build"}}
	if err := db.Db.Where("name = ?", app.Name).FirstOrCreate(&app).Error; err != nil {
		return nil, err
	}
	return &app, err
}

func getOrCreateIns() (*model.Instance, error) {
	var err error
	app, err := getOrCreateApp()
	if err != nil {
		return nil, err
	}
	ins := model.Instance{}
	ins.ApplicationID = app.ID
	ins.Env = "PORT=3001"
	if err := db.Db.Where("application_id = ?", ins.ApplicationID).FirstOrCreate(&ins).Error; err != nil {
		return nil, err
	}
	return &ins, err
}
