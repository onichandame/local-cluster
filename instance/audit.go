package instance

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Audit() error {
	var instances []model.Instance
	if err := db.Db.Find(&instances).Error; err != nil {
		return err
	}
	runningIds, err := listLocalInstances()
	if err != nil {
		return err
	}
	// delete unrecorded instances
	unrecordedIds := []string{}
	for _, id := range runningIds {
		exists := false
		for _, ins := range instances {
			if strconv.Itoa(int(ins.ID)) == id {
				exists = true
			}
		}
		if !exists {
			unrecordedIds = append(unrecordedIds, id)
		}
	}
	for _, unrecordedId := range unrecordedIds {
		go os.Remove(filepath.Join(config.ConfigPresets.InstancesDir, unrecordedId))
	}
	// delete finished instances
	finishedIds := []string{}
	for _, ins := range instances {
		id := strconv.Itoa(int(ins.ID))
		if ins.Status.Value == model.FINISHED && utils.Contains(utils.StrSliceToIfSlice(runningIds), id) {
			finishedIds = append(finishedIds, id)
		}
	}
	for _, finishedId := range finishedIds {
		go os.Remove(filepath.Join(config.ConfigPresets.InstancesDir, finishedId))
	}
	// restart
	statuses := model.GetInstanceStatuses(db.Db)
	policies := model.GetRestartPolicies(db.Db)
	for _, ins := range instances {
		_, ok := RunnersMap[ins.ID]
		if !ok {
			if ins.StatusID == statuses[model.RUNNING].ID {
				// restart running instances
				go runContext(&ins)
			} else if ins.StatusID == statuses[model.FAILED].ID {
				// restart failed instances if restart policy is onfailure or always
				if utils.Contains(utils.UintSliceToIfSlice([]uint{policies[model.ALWAYS].ID, policies[model.ONFAILURE].ID}), ins.RestartPolicyID) {
					go runContext(&ins)
				}
			} else if ins.StatusID == statuses[model.FINISHED].ID {
				// restart finished instances if restart policy is always
				if ins.RestartPolicyID == policies[model.ALWAYS].ID {
					go runContext(&ins)
				}
			}
		}
	}
	return nil
}

func listLocalInstances() ([]string, error) {
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
