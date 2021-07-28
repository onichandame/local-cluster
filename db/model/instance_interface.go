package model

import "gorm.io/gorm"

type InstanceInterface struct {
	gorm.Model
	Port         uint `gorm:"not null"`
	DefinitionID uint `gorm:"not null"`
	Definition   *ApplicationInterface
	InstanceID   uint
}
