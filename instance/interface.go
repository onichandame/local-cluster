package instance

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/onichandame/local-cluster/config"
	appConstants "github.com/onichandame/local-cluster/constants/application"
	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

type InterfaceManager struct {
	lock  sync.Mutex
	ports map[uint]interface{}
	last  uint
}

func (m *InterfaceManager) nextAvailablePort() (port uint) {
	maxTrials := config.Config.PortRange.EndAt - config.Config.PortRange.StartAt + 1
	var trials uint
	var getNextCandidate func(uint) uint
	getNextCandidate = func(p uint) (r uint) {
		if trials >= maxTrials {
			panic(errors.New("no available ports left!"))
		}
		trials++
		if p < config.Config.PortRange.StartAt || p > config.Config.PortRange.EndAt {
			r = config.Config.PortRange.StartAt
		} else {
			if _, ok := m.ports[p]; ok {
				return getNextCandidate(r + 1)
			}
			if tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", r)); err != nil {
				return getNextCandidate(r + 1)
			} else {
				defer tcpListener.Close()
			}
			if udpListener, err := net.Listen("udp", fmt.Sprintf(":%d", r)); err != nil {
				return getNextCandidate(r + 1)
			} else {
				defer udpListener.Close()
			}
		}
		return r
	}
	port = getNextCandidate(m.last)
	return port
}

func (m *InterfaceManager) allocate() (port uint, err error) {
	defer utils.RecoverFromError(&err)
	m.lock.Lock()
	defer m.lock.Unlock()
	port = m.nextAvailablePort()
	m.last = port
	return port, err
}

var ifMan *InterfaceManager

func getInterfaceManager() *InterfaceManager {
	if ifMan == nil {
		atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&ifMan)), nil, unsafe.Pointer(&InterfaceManager{ports: make(map[uint]interface{})}))
	}
	return ifMan
}

func auditInsIfs(rawIns *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	var ins model.Instance
	if err = db.Db.Preload("Interfaces").First(&ins, rawIns.ID).Error; err != nil {
		panic(err)
	}
	manager := getInterfaceManager()
	manager.lock.Lock()
	defer manager.lock.Unlock()
	switch ins.Status {
	// should have interfaces
	case insConstants.CREATING, insConstants.RESTARTING, insConstants.RUNNING:
		var app model.Application
		if err = db.Db.Preload("LocalApplication.Interfaces").Preload("RemoteApplication.Interfaces").First(&app, "name = ?", ins.ApplicationName).Error; err != nil {
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
					if insIf.Port, err = manager.allocate(); err != nil {
						panic(err)
					}
					if err = db.Db.Create(insIf).Error; err != nil {
						panic(err)
					}
				} else {
					if _, ok := manager.ports[insIf.Port]; !ok {
						manager.ports[insIf.Port] = nil
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
					}
				}
			}
		case appConstants.STATIC:
			if len(ins.Interfaces) < 1 {
				var insIf model.InstanceInterface
				insIf.InstanceID = ins.ID
				if insIf.Port, err = manager.allocate(); err != nil {
					panic(err)
				}
				if err = db.Db.Create(&insIf).Error; err != nil {
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
	case insConstants.TERMINATING, insConstants.TERMINATED, insConstants.CRASHED:
		for _, insIf := range ins.Interfaces {
			delete(manager.ports, insIf.Port)
			if err = db.Db.Delete(&insIf).Error; err != nil {
				panic(err)
			}
		}
	}
	return err
}
