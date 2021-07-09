package model

import "gorm.io/gorm"

type InstanceTemplate struct {
	ApplicationID   uint
	Application     Application
	RestartPolicyID uint
	RestartPolicy   Enum
	Env             string
	Port            string
	InstanceGroupID uint
}

type InstanceGroup struct {
	gorm.Model
	Replicas uint
	Template InstanceTemplate
}
