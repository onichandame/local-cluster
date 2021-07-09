package instance

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
)

func Audit() error {
	runtimes, err := listRuntimes()
	if err != nil {
		return err
	}
	// delete unrecorded runtimes
	unrecordedRuntimes := []string{}
	for _, runtime := range runtimes {
		exists := false
		if err := db.Db.Where("id = ?", runtime).First(&model.Instance{}).Error; err == nil {
			exists = true
		}
		if !exists {
			unrecordedRuntimes = append(unrecordedRuntimes, runtime)
		}
	}
	for _, unrecordedRuntime := range unrecordedRuntimes {
		logrus.Warnf("deleting unrecorded runtime %s", unrecordedRuntime)
		go os.Remove(filepath.Join(config.ConfigPresets.InstancesDir, unrecordedRuntime))
	}
	// handle recognized runtimes
	for _, runtime := range runtimes {
		id, err := utils.StrToUint(runtime)
		if err != nil {
			logrus.Warnf("cannot parse instance id %s", runtime)
			continue
		}
		ins := model.Instance{}
		if err := db.Db.Where("id = ?", id).First(&ins).Error; err != nil {
			logrus.Warnf("cannot find instance for runtime %d", id)
		}
		switch ins.Status {
		case constants.RUNNING:
			if _, ok := RunnersMap[id]; !ok {
				setInstanceState(&ins, constants.CRASHED)
				switch ins.RestartPolicy {
				case constants.ALWAYS:
					RunInstance(&ins)
				}
			}
		case constants.TERMINATED, constants.CRASHED:
			if _, ok := RunnersMap[id]; ok {
				if err := Terminate(&ins); err != nil {
					logrus.Warnf("instance %d failed to terminate!", id)
					setInstanceState(&ins, constants.CRASHED)
				}
			}
		}
	}
	return nil
}

func listRuntimes() ([]string, error) {
	res := []string{}
	items, err := ioutil.ReadDir(config.ConfigPresets.InstancesDir)
	if err != nil {
		return res, err
	}
	for _, item := range items {
		if item.IsDir() {
			res = append(res, item.Name())
		}
	}
	return res, nil
}
