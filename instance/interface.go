package instance

import (
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/interfaces"
)

func createInterface(insDef *model.Instance) (*model.InstanceInterface, error) {
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
