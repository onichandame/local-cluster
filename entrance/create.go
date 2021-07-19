package entrance

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func Create(ent *model.Entrance) error {
	if ent.ID != 0 {
		return errors.New("record already created")
	}
	return db.Db.Create(ent).Error
}
