package model

import (
	"github.com/onichandame/local-cluster/constants/instance_group"
	"gorm.io/gorm"
)

type InstanceGroup struct {
	gorm.Model
	Selectable
	TemplateID uint `gorm:"not null"`
	Replicas   uint
	Status     instancegroup.InstanceGroupStatus
	Instances  []Instance
}
