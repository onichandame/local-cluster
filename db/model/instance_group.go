package model

import (
	"github.com/onichandame/local-cluster/constants/instance_group"
	"gorm.io/gorm"
)

type InstanceGroup struct {
	gorm.Model
	Selectable
	TemplateName string `gorm:"not null"`
	Replicas     uint
	Status       instancegroup.InstanceGroupStatus
	Instances    []Instance
}
