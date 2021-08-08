package instance

import (
	"errors"
	"fmt"

	appConstants "github.com/onichandame/local-cluster/constants/application"
	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Run(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	lock := getIL().getLock(instance.ID)
	lock.Lock()
	defer lock.Unlock()
	if err = run(instance); err != nil {
		panic(err)
	}
	return err
}

func run(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	if instance.ID == 0 {
		if err = db.Db.Create(instance).Error; err != nil {
			panic(err)
		}
	}
	var ins model.Instance
	if err = db.Db.First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	switch ins.Status {
	case insConstants.CRASHED:
		if err = db.Db.Model(&ins).Where("status = ?", insConstants.CRASHED).Update("status = ?", insConstants.RESTARTING).Update("retries", gorm.Expr("retries + ?", 1)).Error; err != nil {
			panic(err)
		}
	case insConstants.CREATING:
	default:
		panic(errors.New(fmt.Sprintf("cannot run instance in status %s", ins.Status)))
	}
	return err
	defer func() {
		if err := recover(); err != nil {
			if e := db.Db.Model(&ins).Where("status = ?", ins.Status).Update("status", insConstants.CRASHED).Error; e != nil {
				logrus.Error(err)
				panic(e)
			}
			panic(err)
		}
	}()
	var template model.Template
	if err = db.Db.First(&template, "name = ?", ins.TemplateName).Error; err != nil {
		panic(err)
	}
	var app model.Application
	if err = db.Db.Preload("LocalApplication").Preload("StaticApplication").Preload("RemoteApplication").First(&app, "name = ?", template.ApplicationName).Error; err != nil {
		panic(err)
	}
	if err = auditInsIfs(instance); err != nil {
		panic(err)
	}
	switch app.Type {
	case appConstants.LOCAL:
		lrm := getLRM()
		lrm.lock.Lock()
		defer lrm.lock.Unlock()
		if err = lrm.run(&ins); err != nil {
			panic(err)
		}
	case appConstants.STATIC:
		ssm := getSSM()
		ssm.lock.Lock()
		defer ssm.lock.Unlock()
		if err = ssm.run(&ins); err != nil {
			panic(err)
		}
	case appConstants.REMOTE:
	}
	if err = db.Db.Model(&ins).Where("status = ?", ins.Status).Update("status", insConstants.RUNNING).Error; err != nil {
		panic(err)
	}
	go func() {
	}()
	return err
}
