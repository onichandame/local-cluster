package model

import (
	"github.com/onichandame/local-cluster/constants/instance_group"
	"gorm.io/gorm"
)

type InstanceGroup struct {
	gorm.Model
	Replicas uint
	Status   instancegroup.InstanceGroupStatus
	InstanceTemplate
	Instances []Instance
}
