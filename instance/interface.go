package instance

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/onichandame/local-cluster/config"
	appConstants "github.com/onichandame/local-cluster/constants/application"
	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var last uint64

func createIf(insIf *model.InstanceInterface) (err error) {
	defer utils.RecoverFromError(&err)
	oldLast := last
	maxTrials := config.Config.PortRange.EndAt - config.Config.PortRange.StartAt + 1
	var trials uint
	var try func(uint)
	try = func(p uint) {
		tryNext := func() { try(p + 1) }
		if trials >= maxTrials {
			panic(errors.New("no available ports left!"))
		}
		trials++
		normalize := func() uint {
			if p < config.Config.PortRange.StartAt || p > config.Config.PortRange.EndAt {
				p = config.Config.PortRange.StartAt
			}
			var insIf model.InstanceInterface
			if err := db.Db.First(&insIf, "port = ?", p).Error; err == nil {
				return 0
			} else {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					panic(err)
				}
			}
			if tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", p)); err != nil {
				return 0
			} else {
				defer tcpListener.Close()
			}
			if udpListener, err := net.Listen("udp", fmt.Sprintf(":%d", p)); err != nil {
				return 0
			} else {
				defer udpListener.Close()
			}
			return p
		}
		port := normalize()
		if port == 0 {
			tryNext()
			return
		}
		insIf.Port = port
		if err := db.Db.Create(insIf).Error; err != nil {
			logrus.Warn(err)
			tryNext()
			return
		}
		atomic.CompareAndSwapUint64(&last, oldLast, uint64(insIf.Port))
	}
	try(uint(last + 1))
	return err
}

func auditInsIfs(rawIns *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	var ins model.Instance
	if err := db.Db.Preload("Interfaces").Preload("Template").First(&ins, rawIns.ID).Error; err != nil {
		panic(err)
	}
	switch ins.Status {
	// should have interfaces
	case insConstants.CREATING, insConstants.WAITING, insConstants.RESTARTING, insConstants.RUNNING:
		var app model.Application
		if err := db.Db.Preload("LocalApplication.Interfaces").Preload("RemoteApplication.Interfaces").First(&app, "name = ?", ins.Template.ApplicationName).Error; err != nil {
			panic(err)
		}
		switch app.Type {
		case appConstants.LOCAL:
			for _, appIf := range app.LocalApplication.Interfaces {
				var insIf *model.InstanceInterface
				for _, i := range ins.Interfaces {
					if i.DefinitionName == appIf.Name {
						insIf = &i
					}
				}
				if insIf == nil {
					insIf = &model.InstanceInterface{}
					insIf.InstanceID = ins.ID
					insIf.DefinitionName = appIf.Name
					if err = createIf(insIf); err != nil {
						panic(err)
					}
				}
			}
			for _, insIf := range ins.Interfaces {
				var appIf *model.LocalApplicationInterface
				for _, i := range app.LocalApplication.Interfaces {
					if i.Name == insIf.DefinitionName {
						appIf = &i
					}
				}
				if appIf != nil {
					if err = db.Db.Delete(&insIf).Error; err != nil {
						panic(err)
					}
				}
			}
		case appConstants.STATIC:
			if len(ins.Interfaces) < 1 {
				var insIf model.InstanceInterface
				insIf.InstanceID = ins.ID
				if err = createIf(&insIf); err != nil {
					panic(err)
				}
			} else if len(ins.Interfaces) > 1 {
				willDelete := len(ins.Interfaces) - 1
				deleted := 0
				for _, insIf := range ins.Interfaces {
					if err = db.Db.Delete(&insIf).Error; err != nil {
						panic(err)
					}
					deleted++
					if deleted >= willDelete {
						break
					}
				}
			}
		case appConstants.REMOTE:
			for _, insIf := range ins.Interfaces {
				if err = db.Db.Delete(&insIf).Error; err != nil {
					panic(err)
				}
			}
		}
		// should not have interfaces
	default:
		for _, insIf := range ins.Interfaces {
			if err = db.Db.Delete(&insIf).Error; err != nil {
				panic(err)
			}
		}
	}
	return err
}

func auditOrphanIfs() (err error) {
	defer utils.RecoverFromError(&err)
	if rows, err := db.Db.Model(&model.InstanceInterface{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		var wg sync.WaitGroup
		for rows.Next() {
			wg.Add(1)
			var insIf model.InstanceInterface
			if err := db.Db.ScanRows(rows, &insIf); err != nil {
				panic(err)
			}
			go func() {
				defer wg.Done()
				if err := db.Db.First(&model.Instance{}, insIf.InstanceID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						if err = db.Db.Delete(&insIf).Error; err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}
				}
			}()
		}
		wg.Wait()
	}
	return err
}
