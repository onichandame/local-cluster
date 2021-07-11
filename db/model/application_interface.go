package model

import "gorm.io/gorm"

type ApplicationInterface struct {
	gorm.Model
	ApplicationID uint `gorm:"not null"`
	Name          string
	PortByEnv     string
	PortByArg     string
}
