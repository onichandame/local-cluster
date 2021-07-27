package model

import "gorm.io/gorm"

type Gateway struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Port        uint   `gorm:"not null"`
	ServiceID   uint
	Service     Service
	InterfaceID uint
	Interface   ApplicationInterface
}
