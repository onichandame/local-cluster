package model

import (
	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type JobRecord struct {
	gorm.Model
	Job      string `gorm:"index:linked_list,unique"`
	StatusID uint
	Status   constants.JobStatus
	PrevID   uint `gorm:"index:linked_list,unique"`
	Prev     *JobRecord
	Output   string
	RunByID  uint
	RunBy    User
}
