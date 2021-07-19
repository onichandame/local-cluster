package instance

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

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
	// audit not-settled instances
	notSettledInstances := []model.Instance{}
	if err := db.Db.Where("status NOT IN ?", []constants.InstanceStatus{constants.TERMINATED, constants.CRASHED}).Find(&notSettledInstances).Error; err != nil {
		return err
	}
	for _, ins := range notSettledInstances {
		_, running := RunnersMap[ins.ID]
		if !utils.Contains(utils.StrSliceToIfSlice(runtimes), strconv.Itoa(int(ins.ID))) || !running {
			setInstanceState(&ins, constants.CRASHED)
		}
	}
	// remove unused runtimes
	for _, runtime := range runtimes {
		id, err := utils.StrToUint(runtime)
		if err != nil {
			logrus.Warnf("cannot parse instance id %s", runtime)
			continue
		}
		ins := model.Instance{}
		if err := db.Db.Where("id = ? AND status IN ?", id, []constants.InstanceStatus{constants.CREATING, constants.RUNNING, constants.TERMINATING}).First(&ins).Error; err != nil {
			logrus.Warnf("cannot find active instance for runtime %d. deleting it", id)
			if err := os.RemoveAll(filepath.Join(config.Config.Path.Instances, runtime)); err != nil {
				logrus.Warnf("failed to delete unused runtime %d", id)
			}
		}
	}
	// restart crashed instance if required
	crashedRows, err := db.Db.Model(&model.Instance{}).Where("status = ?", constants.CRASHED).Rows()
	defer crashedRows.Close()
	var crashedRow model.Instance
	for crashedRows.Next() {
		if err := db.Db.ScanRows(crashedRows, &crashedRow); err != nil {
			return err
		}
		switch crashedRow.RestartPolicy {
		case constants.ALWAYS:
			RunInstance(&crashedRow)
		}
	}
	return nil
}

func listRuntimes() ([]string, error) {
	res := []string{}
	items, err := ioutil.ReadDir(config.Config.Path.Instances)
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
