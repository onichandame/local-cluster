package instancegroup

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type InstanceGroupLocks struct {
	lock   sync.Mutex
	groups map[uint]*sync.Mutex
}

func (il *InstanceGroupLocks) getLock(insId uint) *sync.Mutex {
	il.lock.Lock()
	defer il.lock.Unlock()
	if _, ok := il.groups[insId]; !ok {
		il.groups[insId] = &sync.Mutex{}
	}
	return il.groups[insId]
}

var il *InstanceGroupLocks

func getIGL() *InstanceGroupLocks {
	if il == nil {
		atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&il)), nil, unsafe.Pointer(&InstanceGroupLocks{groups: make(map[uint]*sync.Mutex)}))
	}
	return il
}
