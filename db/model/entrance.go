package model

import "gorm.io/gorm"

type Entrance struct {
	gorm.Model
	InstanceID uint
	Instance   Instance
	Path       string
}
