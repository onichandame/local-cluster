package model

import "gorm.io/gorm"

type InstanceInterface struct {
	gorm.Model
	InstanceID uint `gorm:"not null"`
	Port       uint
}
