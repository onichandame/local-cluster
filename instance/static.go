package instance

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/gin-gonic/gin"
)

type StaticServerManager struct {
	lock    sync.Mutex
	servers []*gin.Engine
}

var ssm *StaticServerManager

func getStaticServerManager() *StaticServerManager {
	if ssm == nil {
		atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&ssm)), nil, unsafe.Pointer(&StaticServerManager{servers: make([]*gin.Engine, 0)}))
	}
	return ssm
}
