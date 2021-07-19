package model

import "gorm.io/gorm"

type ServiceInterface struct {
	gorm.Model
	Port         uint `gorm:"not null"`
	DefinitionID uint `gorm:"not null"`
	Definition   *ApplicationInterface
	ServiceID    uint
	ServiceType  string
}
