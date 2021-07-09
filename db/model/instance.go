package model

import (
	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type Instance struct {
	gorm.Model
	InstanceTemplate
	StatusID uint
	Status   constants.InstanceStatus
	GroupID  uint
}
