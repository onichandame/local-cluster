package model

import (
	jobConstants "github.com/onichandame/local-cluster/constants/job"
	"gorm.io/gorm"
)

type JobRecord struct {
	gorm.Model
	Job     string
	Status  jobConstants.JobStatus
	PrevID  uint `gorm:"unique"`
	Prev    *JobRecord
	Output  string
	RunByID uint
	RunBy   User
}
