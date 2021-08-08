package job

import (
	"errors"
	"fmt"
	"time"

	jobConstants "github.com/onichandame/local-cluster/constants/job"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"gorm.io/gorm"
)

var errJobAlreadyRun = errors.New(fmt.Sprintf("job already created by another runner"))

func runJob(j *job) (err error) {
	defer func() {
		if err != nil {
			if j.fatal {
				panic(err)
			}
		}
	}()
	defer utils.RecoverFromError(&err)
	if j.totalRuns != 0 {
		var count int64
		if err = db.Db.Model(&model.JobRecord{}).Where("job = ?", j.name).Count(&count).Error; err != nil {
			panic(err)
		}
		if count >= int64(j.totalRuns) {
			panic(errors.New(fmt.Sprintf("job %s has depleted all allowed runs", j.name)))
		}
	}
	if j.successfulRuns != 0 {
		var count int64
		if err = db.Db.Model(&model.JobRecord{}).Where("job = ? AND status = ?", j.name, jobConstants.FINISHED).Count(&count).Error; err != nil {
			panic(err)
		}
		if count >= int64(j.successfulRuns) {
			panic(errors.New(fmt.Sprintf("job %s has depleted all allowed successful runs", j.name)))
		}
	}
	var prev model.JobRecord
	if err = db.Db.Order("updated_at desc").First(&prev).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
	}
	if j.interval != 0 {
		if time.Now().Sub(prev.UpdatedAt) < j.interval-time.Second {
			panic(nil)
		}
	}
	var record model.JobRecord
	record.PrevID = prev.ID
	record.Job = j.name
	record.Status = jobConstants.PENDING
	if err = db.Db.Create(&record).Error; err == nil {
		if err = j.run(); err == nil {
			record.Status = jobConstants.FINISHED
		} else {
			record.Status = jobConstants.FAILED
			record.Output = err.Error()
		}
		if err = db.Db.Save(&record).Error; err != nil {
			panic(err)
		}
		if j.fatal {
			if record.Status != jobConstants.FINISHED {
				panic(errors.New(record.Output))
			}
		}
	}
	return err
}
