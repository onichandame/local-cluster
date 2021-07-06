package model

import "gorm.io/gorm"

type Instance struct {
	gorm.Model
	Application   Application
	ApplicationID uint
	StatusID      uint
	Status        Enum
}
