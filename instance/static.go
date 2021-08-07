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
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

type StaticServerManager struct {
	lock    sync.Mutex
	servers map[uint]*http.Server
}

func (m *StaticServerManager) run(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	var ins model.Instance
	if err = db.Db.Preload("Interfaces").First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	if len(ins.Interfaces) < 1 {
		panic(errors.New(fmt.Sprintf("instance %d does not have enough interface prepared!", ins.ID)))
	} else if len(ins.Interfaces) > 1 {
		panic(errors.New(fmt.Sprintf("static app only needs 1 interface! audit instance %d's interfaces then retry", ins.ID)))
	}
	if _, ok := m.servers[ins.ID]; ok {
		panic(errors.New(fmt.Sprintf("instance %d already runned! audit first!")))
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

func (m *StaticServerManager) stop(id uint) (err error) {
	defer utils.RecoverFromError(&err)
	if server, ok := m.servers[id]; ok {
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
