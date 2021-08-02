package model

import "gorm.io/gorm"

type ApplicationInterface struct {
	gorm.Model
	Selectable
	ApplicationID uint `gorm:"not null"`
	PortByEnv     string
	PortByArg     string
}
