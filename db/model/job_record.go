package model

import (
	"github.com/onichandame/local-cluster/constants"
	"gorm.io/gorm"
)

type JobRecord struct {
	gorm.Model
	Job      string
	StatusID uint
	Status   constants.JobStatus
	PrevID   uint `gorm:"unique"`
	Prev     *JobRecord
	Output   string
	RunByID  uint
	RunBy    User
}
