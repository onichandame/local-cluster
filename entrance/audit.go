package entrance

import (
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/sirupsen/logrus"
)

func Audit() error {
	var err error
	entrances := []model.Entrance{}
	if err = db.Db.Preload("Backend").Find(&entrances).Error; err != nil {
		return err
	}
	count := 0
	for _, ent := range entrances {
		if ent.Backend.ID == 0 {
			count++
			if err = db.Db.Delete(&ent).Error; err != nil {
				return err
			}
		}
	}
	logrus.Infof("deleted %d orphan entrances", count)
	return err
}
