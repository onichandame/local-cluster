package model

import "gorm.io/gorm"

type Credential struct {
	gorm.Model
	UserID uint
}
