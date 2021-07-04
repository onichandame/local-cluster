package model

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Status string

const (
	PENDING  Status = "PENDING"
	FINISHED Status = "FINISHED"
	FAILED   Status = "FAILED"
)

type JobStatus struct {
	gorm.Model
	Status Status `gorm:"unique"`
}

func GetJobStatuses(db *gorm.DB) map[Status]*JobStatus {
	recChan := make(chan *JobStatus)
	getRec := func(status Status) {
		rec := JobStatus{}
		err := db.FirstOrCreate(&rec, JobStatus{Status: status}).Error
		if err != nil {
			logrus.Fatalf("failed to get job status records")
		}
		recChan <- &rec
	}
	statuses := []Status{PENDING, FINISHED, FAILED}
	for _, s := range statuses {
		go getRec(s)
	}
	res := make(map[Status]*JobStatus)
	for range statuses {
		rec := <-recChan
		res[rec.Status] = rec
	}
	return res
}
