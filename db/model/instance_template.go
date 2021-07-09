package model

import (
	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type InstanceTemplate struct {
	ApplicationID   uint
	Application     Application
	RestartPolicy   constants.InstanceRestartPolicy
	Env             string
	Port            string
	InstanceGroupID uint
}

type InstanceGroup struct {
	gorm.Model
	Replicas uint
	Template InstanceTemplate
}
