package model

import (
	"github.com/onichandame/local-cluster/constants"
)

type InstanceTemplate struct {
	Selectable
	ApplicationID uint `gorm:"not null"`
	Application   Application
	RestartPolicy constants.InstanceRestartPolicy
	Env           string
}
