package interfaces

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func createIF(insDef *model.Instance, ifDef *model.ApplicationInterface) (err error) {
	defer utils.RecoverFromError(&err)
	if insDef.ApplicationID != ifDef.ApplicationID {
		panic(errors.New("instance and interface definition must belong to the same application!"))
	}
	insDef.Interfaces = append(insDef.Interfaces, model.InstanceInterface{Definition: ifDef})
	if err = db.Db.Save(&insDef).Error; err != nil {
		panic(err)
	}
	if err = register(&insDef.Interfaces[len(insDef.Interfaces)-1]); err != nil {
		panic(err)
	}
	return err
}
