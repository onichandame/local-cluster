package model

import (
	"gorm.io/gorm"
)

type ApplicationSpec struct {
	gorm.Model
	ApplicationID uint
	Platform      string
	Arch          string
	Entrypoint    string
	DownloadUrl   string
	Hash          string
}
