package model

import (
	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type InstanceTemplate struct {
	gorm.Model
	InstanceGroupID uint `gorm:"not null"`
	ApplicationID   uint `gorm:"not null"`
	Application     Application
	RestartPolicy   constants.InstanceRestartPolicy
	Env             string
}
