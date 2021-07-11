package model

import (
	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type InstanceTemplate struct {
	gorm.Model
	InstanceGroupID uint
	ApplicationID   uint
	Application     Application
	RestartPolicy   constants.InstanceRestartPolicy
	Env             string
}
