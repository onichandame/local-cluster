package model

import "gorm.io/gorm"

type ApplicationInterface struct {
	gorm.Model
	ApplicationID uint
	Name          string `gorm:"unique;not null"`
	PortByEnv     string
	PortByArg     string
}
