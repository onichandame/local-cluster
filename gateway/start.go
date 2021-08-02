package gateway

import (
	"errors"
	"time"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/onichandame/local-cluster/proxy"
)

func Start(gwDef *model.Gateway) (err error) {
	defer utils.RecoverFromError(&err)
	if gwDef.ID == 0 {
		if err = db.Db.Create(gwDef).Error; err != nil {
			panic(err)
		}
	}
	if gwDef.Port == 0 {
		panic(errors.New("gateway must have a definded port!"))
	}

	if err = db.Db.Create(gwDef).Error; err != nil {
		panic(err)
	}

	audit := func() (err error) {
		defer utils.RecoverFromError(&err)
		if err = db.Db.First(gwDef, gwDef.ID).Error; err != nil {
			panic(err)
		}
		var instances []model.Instance
		if raws, err := getServices(gwDef); err != nil {
			panic(err)
		} else {
			ids := []uint{}
			for _, raw := range raws {
				ids = append(ids, raw.ID)
			}
			if err = db.Db.Preload("Interfaces.Definition").Where("id IN ?", ids).Find(&instances).Error; err != nil {
				panic(err)
			}
		}
		targets := []uint{}
		for _, ins := range instances {
			for _, insIf := range ins.Interfaces {
				if insIf.Definition.Name == gwDef.InterfaceName {
					targets = append(targets, insIf.Port)
				}
			}
		}
		p := proxy.Read(gwDef.Port)
		restart := func() {
			proxy.Delete(gwDef.Port)
			if err := proxy.Create(gwDef.Port, targets); err != nil {
				panic(err)
			}
		}
		if p == nil {
			restart()
		} else {
			hasTargetsChanged := func() (res bool) {
				if res = len(targets) == len(p.Targets); res {
					return res
				}
				for _, port := range targets {
					if !proxy.HasPort(p, port) {
						res = true
						break
					}
				}
				return res
			}
			if hasTargetsChanged() {
				restart()
			}
		}
		return err
	}

	go func() {
		for {
			if err := audit(); err != nil {
				break
			} else {
				time.Sleep(time.Duration(5) * time.Second)
			}
		}
	}()

	return err
}
