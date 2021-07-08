package instance

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func AuditApps(appDef *model.Application) error {
	var instances []model.Instance
	if err := db.Db.Where("application_id = ?", appDef.ID).Find(&instances).Error; err != nil {
		return err
	}
	runningIds, err := listLocalInstances(appDef)
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
		go os.Remove(filepath.Join(getAppDir(appDef), unrecordedId))
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
		go os.Remove(filepath.Join(getAppDir(appDef), finishedId))
	}
	return nil
}

func listLocalInstances(appDef *model.Application) ([]string, error) {
	res := []string{}
	appDir := getAppDir(appDef)
	items, err := ioutil.ReadDir(appDir)
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
