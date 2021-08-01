package interfaces

import (
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Release(insDef *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	if err := db.Db.Preload("Interfaces").Find(&insDef).Error; err != nil {
		panic(err)
	}
	for _, insIf := range insDef.Interfaces {
		if e := portAllocator.Deallocate(insIf.Port); e != nil {
			err = e
		}
	}
	if err != nil {
		panic(err)
	}
	return err
}
