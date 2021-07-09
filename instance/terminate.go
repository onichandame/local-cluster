package instance

import (
	"errors"

	"github.com/onichandame/local-cluster/db/model"
)

func Terminate(insDef *model.Instance) error {
	runner, ok := RunnersMap[insDef.ID]
	if !ok {
		return errors.New("cannot terminate a not-running instance. crash it!")
	}
	runner.cancel()
	return nil
}
