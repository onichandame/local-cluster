package job

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var errJobAlreadyRun = errors.New(fmt.Sprintf("job already created by another runner"))

func runJob(j *job) error {
	if j.totalRuns != 0 {
		if runs, err := countRuns(struct {
			job     string
			success bool
		}{job: j.name}); err != nil {
			return errors.New(fmt.Sprintf("failed to count job records for job %s", j.name))
		} else if runs >= j.totalRuns {
			return errors.New(fmt.Sprintf("job %s has depleted all allowed runs", j.name))
		}
	}
	if j.successfulRuns != 0 {
		if runs, err := countRuns(struct {
			job     string
			success bool
		}{job: j.name, success: true}); err != nil {
			return errors.New(fmt.Sprintf("failed to count successful runs for job %s", j.name))
		} else if runs >= j.successfulRuns {
			return errors.New(fmt.Sprintf("job %s has depleted all allowed successful runs", j.name))
		}
	}
	prev := findLastRun(j)
	run, err := initiateRun(j, prev)
	if err != nil {
		if !errors.Is(err, errJobAlreadyRun) {
			if j.fatal {
				return errors.New(fmt.Sprintf("job %s failed and it is fatal", j.name))
			}
		}
	}
	if err = j.run(); err == nil {
		return finalizeRun(run, constants.FINISHED)
	} else {
		logrus.Warnf("job %d failed", run.ID)
		logrus.Warn(err)
		finalizeRun(run, constants.FAILED)
		if j.fatal {
			panic(err)
		} else {
			return err
		}
	}
}

func findLastRun(job *job) uint {
	prev := model.JobRecord{}
	err := db.Db.Order("created_at desc").First(&prev).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Errorf("db error: failed to find job record!")
			logrus.Fatal(err)
		}
	}
	return prev.ID
}

func initiateRun(job *job, prev uint) (*model.JobRecord, error) {
	run := model.JobRecord{Job: job.name, Status: constants.PENDING, PrevID: prev}
	err := db.Db.Create(&run).Error
	if err != nil {
		logrus.Error(err)
		err = db.Db.Where("prev_id = ?", "prev").First(&run).Error
		if err == nil {
			return nil, errors.New(fmt.Sprintf("failed to create new record for job %s", job.name))
		} else {
			return nil, errJobAlreadyRun
		}
	}
	logrus.Infof("starting job %s", job.name)
	return &run, nil
}

func finalizeRun(run *model.JobRecord, status constants.JobStatus) error {
	return db.Db.Model(run).Where("status = ?", constants.PENDING).Update("status", status).Error
}

func countRuns(args struct {
	job     string
	success bool
}) (uint, error) {
	var count int64
	query := db.Db.Model(&model.JobRecord{}).Where("job = ?", args.job)
	if args.success {
		query = query.Where("status = ?", constants.FINISHED)
	}
	err := query.Count(&count).Error
	return uint(count), err
}
