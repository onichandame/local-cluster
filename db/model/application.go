package model

import (
	"github.com/onichandame/local-cluster/constants/application"
	"gorm.io/gorm"
)

type downloadable struct {
	DownloadUrl string `gorm:"not null"`
	Hash        string
}

type typedApp struct {
	ApplicationID uint `gorm:"not null"`
}

type Application struct {
	gorm.Model
	Selectable
	Type              application.ApplicationType `gorm:"not null"`
	LocalApplication  *LocalApplication
	StaticApplication *StaticApplication
	RemoteApplication *RemoteApplication
}

type LocalApplication struct {
	gorm.Model
	typedApp
	Specs      []LocalApplicationSpec
	Interfaces []LocalApplicationInterface
}

type LocalApplicationSpec struct {
	gorm.Model
	downloadable
	LocalApplicationID uint   `gorm:"not null"`
	Platform           string `gorm:"not null"`
	Arch               string `gorm:"not null"`
	Command            string
	Args               string
}

type LocalApplicationInterface struct {
	gorm.Model
	selectable
	LocalApplicationID uint `gorm:"not null"`
	PortByEnv          string
	PortByArg          string
}

type StaticApplication struct {
	gorm.Model
	typedApp
	downloadable
}

type RemoteApplication struct {
	gorm.Model
	typedApp
	Host       string `gorm:"not null"`
	Interfaces RemoteApplicationInterface
}

type RemoteApplicationInterface struct {
	gorm.Model
	selectable
	RemoteApplicationID uint `gorm:"not null"`
	Port                uint `gorm:"not null"`
}
