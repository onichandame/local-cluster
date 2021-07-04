package model

import "gorm.io/gorm"

type JobRecord struct {
	gorm.Model
	Job      string
	StatusID uint
	Status   JobStatus
	PrevID   uint `gorm:"unique"`
	Prev     *JobRecord
	Output   string
	RunByID  uint
	RunBy    User
}
