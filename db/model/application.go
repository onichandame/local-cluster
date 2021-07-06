package model

import "gorm.io/gorm"

type Application struct {
	gorm.Model
	Name  string
	Specs []ApplicationSpec
}
