package model

import "gorm.io/gorm"

type Entrance struct {
	gorm.Model
	Path string
}
