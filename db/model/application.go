package model

import "gorm.io/gorm"

type Application struct {
	gorm.Model
	Name       string `gorm:"unique"`
	Specs      []ApplicationSpec
	Interfaces []ApplicationInterface
}
