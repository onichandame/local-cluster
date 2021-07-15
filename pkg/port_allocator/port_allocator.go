package portallocator

import "sync"

type PortAllocator struct {
	ports      map[uint]interface{}
	lastUsed   uint
	lock       sync.Mutex
	lowerBound uint
	upperBound uint
}

func New(lowerBound uint, upperBound uint) *PortAllocator {
	var a PortAllocator
	a.ports = make(map[uint]interface{})
	if lowerBound >= upperBound {
		panic("lower bound must be lower than upper bound!")
	}
	a.lowerBound = lowerBound
	a.upperBound = upperBound
	return &a
}
