package model

import "gorm.io/gorm"

type Instance struct {
	gorm.Model
	InstanceTemplate
	StatusID uint
	Status   Enum
	GroupID  uint
}
