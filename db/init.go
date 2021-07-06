package db

import (
	"path/filepath"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func DBInit() {
	dbDir, err := config.GetDbDir()
	if err != nil {
		panic("failed to load db path")
	}
	dbPath := filepath.Join(dbDir, "core.sqlite")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
	Db.AutoMigrate(&model.Enum{})
	Db.AutoMigrate(&model.Application{})
	Db.AutoMigrate(&model.ApplicationSpec{})
	Db.AutoMigrate(&model.Instance{})
}
