package model

import (
	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model
	Selectable

	ApplicationName uint `gorm:"not null"`
	Probes          []Probe
	MaxRetries      uint
	Env             string
	StorageBindings []StorageBinding
}

type Probe struct {
	gorm.Model
	TemplateID   uint   `gorm:"not null"`
	InitialDelay uint   `gorm:"not null"`
	Period       string `gorm:"not null"`
	TCPProbe     TCPProbe
	HTTPProbe    HTTPProbe
}

type probe struct {
	InterfaceName string `gorm:"not null"`
	ProbeID       uint   `gorm:"not null"`
}

type TCPProbe struct {
	gorm.Model
	probe
}

type HTTPProbe struct {
	gorm.Model
	probe
	Path   string               `gorm:"not null"`
	Method constants.HTTPMethod `gorm:"not null"`
	Status uint
}

type StorageBinding struct {
	gorm.Model
	TemplateID  uint   `gorm:"not null"`
	StorageName uint   `gorm:"not null"`
	Path        string `gorm:"not null"`
}
