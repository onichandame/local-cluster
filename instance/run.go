package instance

import (
	"github.com/onichandame/local-cluster/application"
	"github.com/onichandame/local-cluster/db/model"
)

func RunInstance(insDef *model.Instance) error {
	if err := application.PrepareCache(&insDef.Application); err != nil {
		return err
	}
	if err := initInstance(insDef); err != nil {
		return err
	}
	if err := setInstanceState(insDef, model.CREATING); err != nil {
		return nil
	}
	prepareRuntime(insDef)
	if err := runContext(insDef); err != nil {
		return err
	}
	return nil
}
