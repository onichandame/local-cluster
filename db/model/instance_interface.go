package model

import "gorm.io/gorm"

type InstanceInterface struct {
	gorm.Model
	InstanceID             uint `gorm:"not null"`
	ApplicationInterfaceID uint `gorm:"not null"`
	ApplicationInterface   ApplicationInterface
	Port                   uint
}
