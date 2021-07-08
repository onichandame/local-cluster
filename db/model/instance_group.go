package model

import "gorm.io/gorm"

type InstanceGroupTemplate struct {
	gorm.Model
	InstanceTemplate
}

type InstanceGroup struct {
	gorm.Model
	Replicas uint
	Template InstanceGroupTemplate
}
