package db

import (
	"path/filepath"
	"reflect"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func Init() error {
	dbDir := config.Config.Path.DB
	dbPath := filepath.Join(dbDir, "core.sqlite")
	logrus.Infof("opening or creating db at %s", dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return err
	}
	Db = db

	return loadModels()
}

func loadModels() error {
	loadAModel := func(model interface{}) {
		if err := Db.AutoMigrate(model); err != nil {
			logrus.Error(err)
			if t := reflect.TypeOf(model); t.Kind() == reflect.Ptr {
				modelName := "*" + t.Elem().Name()
				logrus.Fatalf("failed to init table %s", modelName)
			}
		}
	}
	// order matters! do not shuffle randomly
	models := []interface{}{
		&model.Application{},
		&model.LocalApplication{},
		&model.LocalApplicationSpec{},
		&model.LocalApplicationInterface{},
		&model.StaticApplication{},
		&model.RemoteApplication{},
		&model.RemoteApplicationInterface{},

		&model.Storage{},

		&model.InstanceProbe{},
		&model.TCPProbe{},
		&model.HTTPProbe{},
		&model.StorageBinding{},

		&model.Instance{},
		&model.InstanceInterface{},

		&model.InstanceGroup{},

		&model.Credential{},

		&model.Entrance{},

		&model.Gateway{},

		&model.JobRecord{},

		&model.User{},
	}
	for _, m := range models {
		loadAModel(m)
	}
	return nil
}
