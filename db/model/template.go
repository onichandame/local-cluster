package model

import (
	"time"

	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model

	ApplicationName uint `gorm:"not null"`
	Probe           *Probe
	MaxRetries      uint
	Env             string
	StorageBindings []StorageBinding

	Instances      []Instance
	InstanceGroups []InstanceGroup
}

type Probe struct {
	gorm.Model
	TemplateID    uint          `gorm:"not null"`
	InitialDelay  time.Duration `gorm:"not null"`
	Interval      time.Duration `gorm:"not null"`
	InterfaceName string        `gorm:"not null"`
	TCPProbe      *TCPProbe
	HTTPProbe     *HTTPProbe
}

type probe struct {
	ProbeID uint `gorm:"not null"`
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
	Status int
}

type StorageBinding struct {
	gorm.Model
	TemplateID  uint   `gorm:"not null"`
	StorageName uint   `gorm:"not null"`
	Path        string `gorm:"not null"`
}
