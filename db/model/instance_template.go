package model

import (
	"github.com/onichandame/local-cluster/constants"
)

type InstanceTemplate struct {
	ApplicationID   uint
	Application     Application
	RestartPolicy   constants.InstanceRestartPolicy
	Env             string
	InstanceGroupID uint
}
