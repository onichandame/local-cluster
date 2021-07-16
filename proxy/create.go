package proxy

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/pkg/proxy"
)

func Create(srcPort uint, tgtPorts []uint) error {
	proxyManager.lock.Lock()
	defer func() { proxyManager.lock.Unlock() }()
	if _, ok := proxyManager.proxies[srcPort]; ok {
		return errors.New(fmt.Sprintf("port %d already in use!", srcPort))
	}
	var err error
	p := proxy.New()
	p.Source = fmt.Sprintf(":%d", srcPort)
	for _, tgt := range tgtPorts {
		p.Targets = append(p.Targets, fmt.Sprintf("localhost:%d", tgt))
	}
	p.Strategy = proxy.ROUNDROBIN
	err = p.Start()
	proxyManager.proxies[srcPort] = p
	return err
}
