package proxy

import (
	"errors"
	"fmt"
)

func Delete(port uint) error {
	proxyManager.lock.Lock()
	defer func() { proxyManager.lock.Unlock() }()
	proxy, ok := proxyManager.proxies[port]
	if !ok {
		return errors.New(fmt.Sprintf("port %d not in use!", port))
	}
	var err error
	err = proxy.Shutdown()
	delete(proxyManager.proxies, port)
	return err
}
