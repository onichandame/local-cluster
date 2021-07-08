package model

import "gorm.io/gorm"

type InstanceTemplate struct {
	ApplicationID   uint
	Application     Application
	RestartPolicyID uint
	RestartPolicy   Enum
	Env             string
	Port            string
}

type Instance struct {
	gorm.Model
	InstanceTemplate
	StatusID uint
	Status   Enum
	GroupID  uint
	Template InstanceGroupTemplate
}
