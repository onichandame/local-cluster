package interfaces

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Register(ifDef *model.ServiceInterface) error {
	// mutex lock
	Lock.Lock()
	defer func() { Lock.Unlock() }()
	if ifDef.Port != 0 {
		return errors.New(fmt.Sprintf("interface already registered to port %d!", ifDef.Port))
	}
	port, err := lockAPort()
	if err != nil {
		return err
	}
	if err := db.Db.Model(ifDef).Update("port", port).Error; err != nil {
		return err
	}
	delete(LockedPortsMap, port)
	LastRegisteredPort = port
	return nil
}

func lockAPort() (uint, error) {
	start := config.Config.PortRange.StartAt
	end := config.Config.PortRange.EndAt
	port := LastRegisteredPort + 1
	// check if registered
	runningInstances := []model.Instance{}
	if err := db.Db.Where("status IN ?", []constants.InstanceStatus{constants.CREATING, constants.RUNNING}).Find(&runningInstances).Error; err != nil {
		return 0, nil
	}
	runningInsIds := []uint{}
	for _, ins := range runningInstances {
		runningInsIds = append(runningInsIds, ins.ID)
	}
	registeredIfs := []model.ServiceInterface{}
	if err := db.Db.Where("service_id IN ?", runningInsIds).Find(&registeredIfs).Error; err != nil {
		return 0, err
	}
	registeredPortsMap := make(map[uint]interface{})
	for _, i := range registeredIfs {
		registeredPortsMap[i.Port] = nil
	}
	tryPort := func(p uint) bool {
		if p < start || p > end {
			return false
		}
		if _, ok := registeredPortsMap[p]; ok {
			return false
		}
		if _, ok := LockedPortsMap[p]; ok {
			return false
		}
		if !utils.IsPortAvailable(p) {
			return false
		}
		return true
	}
	incPort := func(p *uint) {
		if *p < start || *p > end {
			*p = start
		} else {
			*p = *p + 1
		}
	}
	for !tryPort(port) {
		incPort(&port)
	}
	LockedPortsMap[port] = nil
	return port, nil
}
