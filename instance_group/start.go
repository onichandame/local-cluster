package instancegroup

import (
	"github.com/chebyrash/promise"
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/instance"
	"github.com/sirupsen/logrus"
)

func Start(igDef *model.InstanceGroup) error {
	if igDef.ID == 0 {
		if err := db.Db.Create(&igDef).Error; err != nil {
			return err
		}
	}
	logrus.Infof("starting instance group %d", igDef.ID)
	if err := setInstanceGroupStatus(igDef, constants.INITIALIZING); err != nil {
		return err
	}
	oldInsts := []model.Instance{}
	if err := db.Db.Where("instance_group_id = ?", igDef.ID).Find(&oldInsts).Error; err != nil {
		return err
	}
	for _, oldIns := range oldInsts {
		instance.Terminate(&oldIns)
	}
	ps := []*promise.Promise{}
	logrus.Infof("will start %d replicas", igDef.Replicas)
	for i := 0; i < int(igDef.Replicas); i++ {
		logrus.Infof("starting replica %d", i+1)
		ins := model.Instance{}
		ins.ApplicationID = igDef.ApplicationID
		ins.InstanceGroupID = igDef.ID
		ps = append(ps, promise.New(func(resolve func(promise.Any), reject func(error)) {
			if err := instance.RunInstance(&ins); err == nil {
				resolve(ins)
			} else {
				reject(err)
			}
			return
		}))
	}
	p := promise.All(ps...)
	go func() {
		if _, err := p.Await(); err == nil {
			logrus.Infof("%d replicas started", igDef.Replicas)
			setInstanceGroupStatus(igDef, constants.NOTREADY)
		} else {
			logrus.Error("failed to start replicas")
			logrus.Error(err)
			setInstanceGroupStatus(igDef, constants.READY)
		}
	}()
	return nil
}
