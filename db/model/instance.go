package model

import (
	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type Instance struct {
	gorm.Model
	InstanceTemplate
	Status  constants.InstanceStatus
	GroupID uint
}
