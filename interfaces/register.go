package interfaces

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func register(ifDef *model.ServiceInterface) error {
	if ifDef.Port != 0 {
		return errors.New(fmt.Sprintf("interface already registered to port %d!", ifDef.Port))
	}
	port, err := portAllocator.Allocate()
	if err != nil {
		return err
	}
	if err := db.Db.Model(ifDef).Update("port", port).Error; err != nil {
		return err
	}
	return nil
}
