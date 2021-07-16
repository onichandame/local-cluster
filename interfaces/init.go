package interfaces

import (
	"errors"

	portallocator "github.com/onichandame/local-cluster/pkg/port_allocator"
)

var portAllocator *portallocator.PortAllocator

func Init(l uint, u uint) error {
	if portAllocator != nil {
		return errors.New("cannot init port allocator twice!")
	}
	portAllocator = portallocator.New(l, u)
	return nil
}
