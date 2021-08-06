package model

import (
	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type InstanceTemplate struct {
	Selectable
	ApplicationName uint `gorm:"not null"`
	Probes          []InstanceProbe
	MaxRetries      uint
	Env             string
	StorageBindings []StorageBinding
}

type InstanceProbe struct {
	gorm.Model
	InitialDelay uint   `gorm:"not null"`
	Period       string `gorm:"not null"`
	TCPProbe     TCPProbe
	HTTPProbe    HTTPProbe
}

type probe struct {
	InterfaceName string `gorm:"not null"`
}

type TCPProbe struct {
	gorm.Model
	probe
}

type HTTPProbe struct {
	gorm.Model
	probe
	Path   string
	Method constants.HTTPMethod `gorm:"not null"`
	Status uint
}

type StorageBinding struct {
	gorm.Model
	Selectable
	StorageID  uint `gorm:"not null"`
	Storage    Storage
	InstanceID uint `gorm:"not null"`
	Instance
	Path string `gorm:"not null"`
}
