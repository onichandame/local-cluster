package model

import "gorm.io/gorm"

type Entrance struct {
	gorm.Model
	Root      string `gorm:"unique"`
	BackendID uint
	Backend   ServiceInterface
}
