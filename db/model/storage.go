package model

import "gorm.io/gorm"

type Storage struct {
	gorm.Model
	Selectable
	Validated bool
}
