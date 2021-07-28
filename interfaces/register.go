package interfaces

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func register(ifDef *model.InstanceInterface) (err error) {
	defer utils.RecoverFromError(&err)
	if ifDef.Port != 0 {
		panic(errors.New(fmt.Sprintf("interface already registered to port %d!", ifDef.Port)))
	}
	if port, err := portAllocator.Allocate(); err != nil {
		panic(err)
	} else {
		if err = db.Db.Model(ifDef).Update("port", port).Error; err != nil {
			panic(err)
		}
	}
	return err
}
