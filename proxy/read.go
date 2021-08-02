package proxy

import "github.com/onichandame/local-cluster/pkg/proxy"

func Read(source uint) (proxy *proxy.Proxy) {
	proxyManager.lock.Lock()
	defer proxyManager.lock.Unlock()
	proxy = proxyManager.proxies[source]
	return proxy
}
