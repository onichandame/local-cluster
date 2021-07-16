package proxy

import (
	"errors"
	"sync"

	"github.com/onichandame/local-cluster/pkg/proxy"
)

type ProxyManager struct {
	lock    sync.Mutex
	proxies map[uint]*proxy.Proxy
}

func newProxyManager() *ProxyManager {
	var m ProxyManager
	m.proxies = make(map[uint]*proxy.Proxy)
	return &m
}

var proxyManager *ProxyManager

func Init() error {
	if proxyManager != nil {
		return errors.New("cannot init proxy manager twice")
	}
	proxyManager = newProxyManager()
	return nil
}
