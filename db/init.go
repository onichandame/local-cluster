package db

import (
	"github.com/onichandame/local-cluster/db/model"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func DBInit() {
	dbPathKey := "database"
	viper.SetDefault(dbPathKey, "cluster.sqlite")
	dbPath := viper.GetString(dbPathKey)
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
	Db.AutoMigrate(&model.JobStatus{})
}
