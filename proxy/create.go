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
		p.Targets = append(p.Targets, parseTarget(tgt))
	}
	p.Strategy = proxy.ROUNDROBIN
	err = p.Start()
	proxyManager.proxies[srcPort] = p
	return err
}

func HasPort(p *proxy.Proxy, port uint) bool {
	proxyManager.lock.Lock()
	defer proxyManager.lock.Unlock()
	target := parseTarget(port)
	return p.HasTarget(target)
}

func parseTarget(port uint) string {
	return fmt.Sprintf("localhost:%d", port)
}
