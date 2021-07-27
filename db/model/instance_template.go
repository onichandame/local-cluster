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
	ServiceID       uint `gorm:"not null"`
	Service         Service
	RestartPolicy   constants.InstanceRestartPolicy
	Env             string
}
