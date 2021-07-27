package model

import "gorm.io/gorm"

type Entrance struct {
	gorm.Model
	Name      string `gorm:"unique"`
	BackendID uint
	Backend   Gateway
}
