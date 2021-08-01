package gateway

import (
	"container/list"
	"fmt"

	"github.com/chebyrash/promise"
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func getServices(gwDef *model.Gateway) (err error, instances []*model.Instance) {
	defer utils.RecoverFromError(&err)
	ig := new(model.InstanceGroup)
	ins := new(model.Instance)
	if _, err := promise.All([]*promise.Promise{
		promise.New(func(resolve func(promise.Any), reject func(error)) {
			if err = db.Db.Preload("Instances").Where("name = ?", gwDef.ServiceName).First(ig).Error; err != nil {
				panic(err)
			} else {
				resolve(nil)
			}
		}),
		promise.New(func(resolve func(promise.Any), reject func(error)) {
			if err = db.Db.Where("name = ?", gwDef.ServiceName).First(ins).Error; err != nil {
				panic(err)
			} else {
				resolve(nil)
			}
		}),
	}...,
	).Await(); err != nil {
		panic(err)
	}
	services := list.New()
	if ig != nil {
		for _, ins := range ig.Instances {
			services.PushBack(&ins)
		}
	} else if ins != nil {
		services.PushBack(ins)
	} else {
		panic(fmt.Sprintf("gateway %d has no services available", gwDef.ID))
	}
	serivce := services.Front()
	for {
		if serivce == nil {
			break
		}
		if ins, ok := serivce.Value.(*model.Instance); ok {
			if ins.Status != constants.RUNNING {
				services.Remove(serivce)
			}
		}
		serivce = serivce.Next()
	}
	if services.Len() < 1 {
		panic(fmt.Sprintf("gateway %d does not have running service available"))
	}
	instances = make([]*model.Instance, 0)
	service := services.Front()
	for {
		if serivce == nil {
			break
		}
		if ins, ok := service.Value.(*model.Instance); ok {
			instances = append(instances, ins)
		}
		service = service.Next()
	}
	return err, instances
}
