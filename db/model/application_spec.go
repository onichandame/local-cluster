package model

import (
	"gorm.io/gorm"
)

type ApplicationSpec struct {
	gorm.Model
	ApplicationID uint
	Platform      string
	Arch          string
	Command       string
	DownloadUrl   string
	Args          string
	Hash          string
}
