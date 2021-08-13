package instance

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

type Probe struct {
	stop func()
}

type ProbeManager struct {
	lock   sync.Mutex
	probes map[uint]*Probe
}

func (m *ProbeManager) has(instance *model.Instance) (has bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	has = m._has(instance)
	return has
}

func (m *ProbeManager) _has(instance *model.Instance) (has bool) {
	_, has = m.probes[instance.ID]
	return has
}

func (m *ProbeManager) add(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	m.lock.Lock()
	defer m.lock.Unlock()
	err = m._add(instance)
	return err
}

func (m *ProbeManager) _add(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)

	var ins model.Instance
	if err := db.Db.Preload("Template.Probe").First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	if ins.Template == nil || ins.Template.Probe == nil {
		panic(errors.New(fmt.Sprintf("instance %d has no probe set!", instance.ID)))
	}

	var probe Probe
	stopChan := make(chan interface{})
	probe.stop = func() {
		go func() {
			var err error
			defer utils.RecoverFromError(&err)
			stopChan <- nil
		}()
	}
	m.probes[instance.ID] = &probe

	// probes
	startProbe := func() {
		probeOnce := func() {
			defer utils.RecoverAndLog()
			probe := ins.Template.Probe
			var insIf model.InstanceInterface
			if err = db.Db.First(&insIf, "definition_name = ? AND instance_id =?", probe.InterfaceName, ins.ID).Error; err != nil {
				panic(err)
			}
			var running int32 = 1
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if probe.TCPProbe != nil {
					if conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", insIf.Port), insConstants.PROBE_TIMEOUT); err != nil {
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
					if req, err := http.NewRequest(string(probe.HTTPProbe.Method), fmt.Sprintf("http://localhost:%d%s", insIf.Port, probe.HTTPProbe.Path), nil); err != nil {
						panic(err)
					} else {
						client := http.Client{Timeout: insConstants.PROBE_TIMEOUT}
						if res, err := client.Do(req); err != nil {
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
			lock := getIL().getLock(ins.ID)
			lock.Lock()
			defer lock.Unlock()
			id := ins.ID
			var ins model.Instance
			if err := db.Db.First(&ins, id).Error; err != nil {
				panic(err)
			}
			switch ins.Status {
			case insConstants.WAITING, insConstants.RUNNING, insConstants.CRASHED, insConstants.RESTARTING:
				if err = db.Db.Model(&ins).Update("status", status).Error; err != nil {
					panic(err)
				}
			}
		}
		for {
			select {
			case <-stopChan:
				close(stopChan)
				return
			default:
				probeOnce()
			}
			time.Sleep(ins.Template.Probe.Interval)
		}
	}
	go func() {
		time.Sleep(ins.Template.Probe.InitialDelay)
		startProbe()
	}()
	return err
}

func (m *ProbeManager) del(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	m.lock.Lock()
	defer m.lock.Unlock()
	err = m._del(instance)
	return err
}

func (m *ProbeManager) _del(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	delete(m.probes, instance.ID)
	return err
}

var _pm *ProbeManager

func getPM() *ProbeManager {
	if _pm == nil {
		atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&_pm)), nil, unsafe.Pointer(&ProbeManager{}))
	}
	return _pm
}
