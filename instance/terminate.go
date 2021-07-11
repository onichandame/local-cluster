package instance

import (
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db/model"
)

func Terminate(insDef *model.Instance) error {
	if err := setInstanceState(insDef, constants.TERMINATING); err != nil {
		return err
	}
	if runner, ok := RunnersMap[insDef.ID]; ok {
		runner.cancel()
		runner.cmd.Wait()
	}
	setInstanceState(insDef, constants.TERMINATED)
	return nil
}
