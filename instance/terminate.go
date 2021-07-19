package instance

import (
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/interfaces"
	"github.com/onichandame/local-cluster/proxy"
)

func Terminate(insDef *model.Instance) error {
	if err := setInstanceState(insDef, constants.TERMINATING); err != nil {
		return err
	}
	if runner, ok := RunnersMap[insDef.ID]; ok {
		runner.cancel()
		runner.cmd.Wait()
	}
	if err := db.Db.Preload("Interfaces").First(insDef, insDef.ID).Error; err != nil {
		return err
	}
	for _, insIf := range insDef.Interfaces {
		if err := proxy.Delete(insIf.Port); err != nil {
			return err
		}
	}
	if err := interfaces.ReleaseIF(insDef); err != nil {
		return err
	}
	setInstanceState(insDef, constants.TERMINATED)
	return nil
}
