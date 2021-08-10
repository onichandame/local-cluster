package instance

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/onichandame/local-cluster/constants"
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
	if err = db.Db.Preload("Template").First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	if ins.Template == nil {
		panic(errors.New(fmt.Sprintf("instance %d has invalid template!", ins.ID)))
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
	defer func() {
		if err := recover(); err != nil {
			if e := db.Db.Model(&ins).Where("status = ?", ins.Status).Update("status", insConstants.CRASHED).Error; e != nil {
				logrus.Error(err)
				panic(e)
			}
			panic(err)
		}
	}()
	var app model.Application
	if err = db.Db.Preload("LocalApplication").Preload("StaticApplication").Preload("RemoteApplication").First(&app, "name = ?", ins.Template.ApplicationName).Error; err != nil {
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
	// probes
	go func() {
		for _, probe := range ins.Template.Probes {
			probe := probe
			var insIf model.InstanceInterface
			if err = db.Db.First(&insIf, "definition_name = ? AND instance_id =?", probe.InterfaceName, ins.ID).Error; err != nil {
				panic(err)
			}
			var run func()
			run = func() {
				var running int32 = 1
				var wg sync.WaitGroup
				wg.Add(1)
				go func() {
					defer wg.Done()
					if probe.TCPProbe != nil {
						if conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", insIf.Port)); err != nil {
							atomic.CompareAndSwapInt32(&running, 1, 0)
						} else {
							defer conn.Close()
						}
					}
				}()
				wg.Add(1)
				go func() {
					var err error
					defer func() {
						if err != nil {
							atomic.CompareAndSwapInt32(&running, 1, 0)
						}
						wg.Done()
					}()
					defer utils.RecoverFromError(&err)
					if probe.HTTPProbe != nil {
						var res *http.Response
						switch probe.HTTPProbe.Method {
						case constants.GET:
							if res, err = http.Get(fmt.Sprintf("http://localhost:%d%s", insIf.Port, probe.HTTPProbe.Path)); err != nil {
								panic(err)
							} else {
								if probe.HTTPProbe.Status != 0 {
									if res.StatusCode != probe.HTTPProbe.Status {
										panic(errors.New(fmt.Sprintf("expect status %d but received %d", probe.HTTPProbe.Status, res.StatusCode)))
									}
								}
							}
						}
					}
				}()
				wg.Wait()
				var status insConstants.InstanceStatus
				if running == 1 {
					status = insConstants.RUNNING
				} else {
					status = insConstants.CRASHED
				}
				if err = db.Db.Model(&ins).Update("status", status).Error; err != nil {
					panic(err)
				}
				time.Sleep(probe.Period)
				run()
			}
			go func() {
				time.Sleep(probe.InitialDelay)
				run()
			}()
		}
	}()
	return err
}
