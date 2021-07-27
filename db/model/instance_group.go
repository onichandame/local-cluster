package model

import (
	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type InstanceGroup struct {
	gorm.Model
	Replicas uint
	Status   constants.InstanceGroupStatus
	InstanceTemplate
	Instances []Instance
}
