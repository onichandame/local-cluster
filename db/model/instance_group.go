package model

import "gorm.io/gorm"

type InstanceGroup struct {
	gorm.Model
	Replicas uint
	Template InstanceTemplate
}
