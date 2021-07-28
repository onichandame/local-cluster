package model

import "gorm.io/gorm"

type Entrance struct {
	gorm.Model
	Selectable
	BackendID uint
	Backend   Gateway
}
