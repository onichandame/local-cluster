package model

import "gorm.io/gorm"

type Storage struct {
	gorm.Model
	Selectable
	Path string `gorm:"not null"`
}
