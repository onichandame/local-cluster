package job

import "github.com/onichandame/local-cluster/proxy"

var initProxyManager = job{
	name:      "InitProxyManager",
	immediate: true,
	blocking:  true,
	fatal:     true,
	run: func() error {
		return proxy.Init()
	},
}
