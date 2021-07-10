package db

import (
	"path/filepath"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func DBInit() {
	dbDir := config.ConfigPresets.DbDir
	dbPath := filepath.Join(dbDir, "core.sqlite")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic("failed to open core database")
	}
	Db = db

	loadModels()
}

func loadModels() {
	Db.AutoMigrate(&model.User{})
	Db.AutoMigrate(&model.Credential{})
	Db.AutoMigrate(&model.JobRecord{})
	Db.AutoMigrate(&model.Application{})
	Db.AutoMigrate(&model.ApplicationSpec{})
	Db.AutoMigrate(&model.Instance{})
	Db.AutoMigrate(&model.InstanceGroup{})
	Db.AutoMigrate(&model.Entrance{})
	Db.AutoMigrate(&model.ApplicationInterface{})
	Db.AutoMigrate(&model.InstanceInterface{})
}
