package job

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var errJobAlreadyRun = errors.New(fmt.Sprintf("job already created by another runner"))

func runJob(job *job) {
	if job.totalRuns != 0 {
		runs, err := countRuns(struct {
			job     string
			success bool
		}{job: job.name})
		if err != nil {
			logrus.Fatalf("failed to count job records for job %s", job.name)
		}
		if runs >= job.totalRuns {
			return
		}
	}
	if job.successfulRuns != 0 {
		runs, err := countRuns(struct {
			job     string
			success bool
		}{job: job.name, success: true})
		if err != nil {
			logrus.Fatalf("failed to count successful runs for job %s", job.name)
		}
		if runs >= job.successfulRuns {
			return
		}
	}
	prev := findLastRun(job)
	runID, err := initiateRun(job, prev)
	if err != nil {
		if !errors.Is(err, errJobAlreadyRun) {
			if job.fatal {
				logrus.Fatalf("job %s failed and it is fatal", job.name)
			}
		}
	}
	err = job.run()
	if err == nil {
		finalizeRun(runID, model.FINISHED)
	} else {
		finalizeRun(runID, model.FAILED)
		if job.fatal {
			logrus.Fatalf("job %s failed", job.name)
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

func initiateRun(job *job, prev uint) (uint, error) {
	statuses := model.GetJobStatuses(db.Db)
	run := model.JobRecord{Job: job.name, StatusID: statuses[model.PENDING].ID, PrevID: prev}
	err := db.Db.Create(&run).Error
	if err != nil {
		logrus.Error(err)
		err = db.Db.Where("prev_id = ?", "prev").First(&run).Error
		if err == nil {
			return 0, errors.New(fmt.Sprintf("failed to create new record for job %s", job.name))
		} else {
			return 0, errJobAlreadyRun
		}
	}
	logrus.Infof("starting job %s", job.name)
	return run.ID, nil
}

func finalizeRun(runID uint, status model.EnumValue) {
	allowedStatuses := map[model.EnumValue]interface{}{model.FINISHED: nil, model.FAILED: nil}
	if _, ok := allowedStatuses[status]; !ok {
		logrus.Fatalf("job record finalized with an unsupported status %s", status)
	}
	statuses := model.GetJobStatuses(db.Db)
	db.Db.Model(&model.JobRecord{}).Where("id = ? AND status_id = ?", runID, statuses[model.PENDING].ID).Update("status_id", statuses[status].ID)

	log := func() {
		run := model.JobRecord{}
		err := db.Db.First(&run, runID).Error
		if err == nil {
			switch status {
			case model.FINISHED:
				logrus.Infof("finished job %s", run.Job)
			case model.FAILED:
				logrus.Errorf("failed job %s", run.Job)
			}
		}
	}
	go log()
}

func countRuns(args struct {
	job     string
	success bool
}) (uint, error) {
	var count int64
	query := db.Db.Model(&model.JobRecord{}).Where("job = ?", args.job)
	if args.success {
		statuses := model.GetJobStatuses(db.Db)
		query = query.Where("status_id = ?", statuses[model.FINISHED].ID)
	}
	err := query.Count(&count).Error
	return uint(count), err
}
