package instance

import (
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/interfaces"
)

func Terminate(insDef *model.Instance) error {
	if err := setInstanceState(insDef, constants.TERMINATING); err != nil {
		return err
	}
	if runner, ok := RunnersMap[insDef.ID]; ok {
		runner.cancel()
		runner.cmd.Wait()
	}
	if err := interfaces.ReleaseIF(insDef); err != nil {
		return err
	}
	setInstanceState(insDef, constants.TERMINATED)
	return nil
}
