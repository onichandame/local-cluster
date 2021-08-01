package instance

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Terminate(insDef *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	manager := getRunnerManager()
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if err = db.Db.First(insDef, insDef.ID).Error; err != nil {
		panic(err)
	}
	if insDef.Status != constants.RUNNING {
		panic(errors.New(fmt.Sprintf("cannot terminate instance in state %s", insDef.Status)))
	}
	if err = setInstanceState(insDef, constants.TERMINATING); err != nil {
		return err
	}
	runner := manager.runners[insDef.ID]
	if runner == nil {
		setInstanceState(insDef, constants.CRASHED)
		panic(errors.New(fmt.Sprintf("instance %d is broken! please run audit", insDef.ID)))
	}
	runner.cancel()
	return err
}
