package model

import "gorm.io/gorm"

type Gateway struct {
	gorm.Model
	Selectable
	Port          uint   `gorm:"not null"`
	ServiceName   string `gorm:"not null"`
	InterfaceName string `gorm:"not null"`
}
