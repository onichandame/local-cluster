package model

import "gorm.io/gorm"

type Application struct {
	gorm.Model
	Selectable
	Specs      []ApplicationSpec
	Interfaces []ApplicationInterface
}
