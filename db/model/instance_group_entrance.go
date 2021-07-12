package model

import "gorm.io/gorm"

type InstanceGroupEntrance struct {
	gorm.Model
	Port            uint
	InstanceGroupID uint `gorm:"not null"`
	InstanceGroup   InstanceGroup
	InterfaceID     uint `gorm:"not null"`
	Interface       ApplicationInterface
}
