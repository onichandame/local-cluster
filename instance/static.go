package instance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

type StaticServerManager struct {
	lock    sync.Mutex
	servers map[uint]*http.Server
}

func (m *StaticServerManager) run(instance *model.Instance) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	err = m._run(instance)
	return err
}

func (m *StaticServerManager) _run(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	var ins model.Instance
	if err := db.Db.Preload("Interfaces").First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	if len(ins.Interfaces) < 1 {
		panic(errors.New(fmt.Sprintf("instance %d does not have enough interface prepared!", ins.ID)))
	} else if len(ins.Interfaces) > 1 {
		panic(errors.New(fmt.Sprintf("static app only needs 1 interface! audit instance %d's interfaces then retry", ins.ID)))
	}
	if _, ok := m.servers[ins.ID]; ok {
		panic(errors.New(fmt.Sprintf("instance %d already runned! audit first!", ins.ID)))
	}
	port := ins.Interfaces[0].Port
	insDir := getInsRuntimeDir(ins.ID)
	handler := gin.New()
	handler.Static("/", insDir)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	m.servers[ins.ID] = &server
	go func() {
		server.ListenAndServe()
		m.lock.Lock()
		defer m.lock.Unlock()
		delete(m.servers, ins.ID)
	}()
	return err
}

func (m *StaticServerManager) stop(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	m.lock.Lock()
	defer m.lock.Unlock()
	err = m._stop(instance)
	return err
}

func (m *StaticServerManager) _stop(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	if server, ok := m.servers[instance.ID]; ok {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = server.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
	return err
}

var ssm *StaticServerManager

func getSSM() *StaticServerManager {
	if ssm == nil {
		atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&ssm)), nil, unsafe.Pointer(&StaticServerManager{servers: make(map[uint]*http.Server)}))
	}
	return ssm
}

func (m *StaticServerManager) audit() (err error) {
	defer utils.RecoverFromError(&err)
	m.lock.Lock()
	defer m.lock.Unlock()
	var wg sync.WaitGroup
	for id := range ssm.servers {
		id := id
		wg.Add(1)
		go func() {
			defer utils.RecoverAndLog()
			defer wg.Done()
			var ins model.Instance
			if err := db.Db.First(&ins, id).Error; err != nil {
				panic(err)
			}
			if ins.Status != insConstants.RUNNING {
				if err := ssm.stop(&ins); err != nil {
					panic(err)
				}
			}
		}()
	}
	wg.Wait()
	return err
}
