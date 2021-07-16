package portallocator

func (a *PortAllocator) Deallocate(port uint) error {
	a.lock.Lock()
	defer func() { a.lock.Unlock() }()
	delete(a.ports, port)
	return nil
}
