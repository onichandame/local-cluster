package model

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Name   string `gorm:"unique,not null"`
	Custom string
}
