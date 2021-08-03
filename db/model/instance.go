package model

import (
	"github.com/onichandame/local-cluster/constants/instance"
	"gorm.io/gorm"
)

type Instance struct {
	gorm.Model
	InstanceTemplate
	Status          instance.InstanceStatus
	InstanceGroupID uint
	Interfaces      []InstanceInterface
	Retries         uint
}

type InstanceInterface struct {
	gorm.Model
	Port           uint   `gorm:"not null"`
	DefinitionName string `gorm:"not null"`
	InstanceID     uint
}
