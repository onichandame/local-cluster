package instance

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/interfaces"
)

func createInterface(insDef *model.Instance, ifDef *model.ApplicationInterface) (*model.InstanceInterface, error) {
	if insDef.ApplicationID != ifDef.ApplicationID {
		return nil, errors.New("interface definition and instance must point to the smae application!")
	}
	insIf := model.InstanceInterface{}
	insIf.InstanceID = insDef.ID
	if err := db.Db.Create(&insIf).Error; err != nil {
		return nil, err
	}
	if err := interfaces.Register(&insIf); err != nil {
		return nil, err
	}
	return &insIf, nil
}
