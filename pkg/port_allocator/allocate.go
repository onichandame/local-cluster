package portallocator

import "errors"

func (a *PortAllocator) nextPort(p uint) uint {
	p++
	if p < a.lowerBound || p >= a.upperBound {
		return a.lowerBound
	} else {
		return p
	}
}

func (a *PortAllocator) Allocate() (uint, error) {
	a.lock.Lock()
	defer func() { a.lock.Unlock() }()
	var port uint = a.lastUsed
	var err error
	tryPort := func(p uint) bool {
		if p < a.lowerBound || p > a.upperBound {
			return false
		}
		if _, ok := a.ports[p]; ok {
			return false
		}
		if err := checkPort(p); err != nil {
			return false
		}
		return true
	}
	var trials uint
	maxTrials := a.upperBound - a.lowerBound + 1
	for !tryPort(port) {
		if trials > maxTrials {
			err = errors.New("all ports in use. cannot allocate more!")
			break
		}
		port = a.nextPort(port)
		trials++
	}
	if err == nil {
		a.ports[port] = nil
		a.lastUsed = port
	}
	return port, err
}
