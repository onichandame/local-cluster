package model

import (
	"github.com/onichandame/local-cluster/constants/instance"
	"gorm.io/gorm"
)

type Instance struct {
	gorm.Model
	Selectable
	TemplateName    string `gorm:"not null"`
	Status          instance.InstanceStatus
	InstanceGroupID uint
	Interfaces      []InstanceInterface
	Retries         uint
}

// Instance creating/restarting → (create)
// Instance crashed/terminating → (delete)

type InstanceInterface struct {
	gorm.Model
	InstanceID     uint `gorm:"not null"`
	Port           uint `gorm:"not null,unique"`
	DefinitionName string
}
