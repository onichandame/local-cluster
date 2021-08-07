package instance

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type InstanceLocks struct {
	lock      sync.Mutex
	instances map[uint]*sync.Mutex
}

func (il *InstanceLocks) getLock(insId uint) *sync.Mutex {
	il.lock.Lock()
	defer il.lock.Unlock()
	if _, ok := il.instances[insId]; !ok {
		il.instances[insId] = &sync.Mutex{}
	}
	return il.instances[insId]
}

var il *InstanceLocks

func getIL() *InstanceLocks {
	if il == nil {
		atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&il)), nil, unsafe.Pointer(&InstanceLocks{instances: make(map[uint]*sync.Mutex)}))
	}
	return il
}
