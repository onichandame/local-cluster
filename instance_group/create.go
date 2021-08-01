package instancegroup

import (
	"errors"
	"fmt"

	"github.com/chebyrash/promise"
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/instance"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
)

func Create(igDef *model.InstanceGroup) (err error) {
	defer utils.RecoverFromError(&err)
	if igDef.ID == 0 {
		igDef.Status = constants.INITIALIZING
		if err = db.Db.Create(&igDef).Error; err != nil {
			panic(err)
		}
	} else {
		panic(errors.New(fmt.Sprintf("instance group %d already created", igDef.ID)))
	}
	switch igDef.Status {
	case constants.INITIALIZING:
	default:
		panic(errors.New(fmt.Sprintf("instance group %d has already initialized", igDef.ID)))
	}
	logrus.Infof("starting instance group %d", igDef.ID)
	defer func() {
		var finalState constants.InstanceGroupStatus
		if err == nil {
			finalState = constants.READY
		} else {
			finalState = constants.NOTREADY
		}
		setInstanceGroupStatus(igDef, finalState)
	}()
	ps := []*promise.Promise{}
	logrus.Infof("will start %d replicas", igDef.Replicas)
	for i := 0; i < int(igDef.Replicas); i++ {
		logrus.Infof("starting replica %d", i+1)
		ins := model.Instance{}
		ins.ApplicationID = igDef.ApplicationID
		ins.InstanceGroupID = igDef.ID
		ins.RestartPolicy = igDef.RestartPolicy
		ins.Env = igDef.Env
		ins.Name = fmt.Sprintf("%s-%d", igDef.Name, i)
		ps = append(ps, promise.New(func(resolve func(promise.Any), reject func(error)) {
			if err = instance.Run(&ins); err == nil {
				resolve(ins)
			} else {
				reject(err)
			}
			return
		}))
	}
	if _, err = promise.All(ps...).Await(); err == nil {
		logrus.Infof("%d replicas started", igDef.Replicas)
	} else {
		logrus.Error("failed to start replicas")
		logrus.Error(err)
		panic(err)
	}
	return err
}
