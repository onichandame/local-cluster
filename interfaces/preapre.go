package interfaces

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func PrepareInterfaces(insDef *model.Instance) (err error) {
	if len(insDef.Interfaces) > 0 {
		panic(errors.New("interfaces already prepared!"))
	}
	if err := db.Db.Preload("Application.Interfaces").First(insDef, insDef.ID).Error; err != nil {
		panic(err)
	}
	for _, ifDef := range insDef.Application.Interfaces {
		if err := createIF(insDef, &ifDef); err != nil {
			panic(err)
		}
	}
	return err
}
