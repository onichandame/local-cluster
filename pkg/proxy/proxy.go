package proxy

import (
	"github.com/onichandame/local-cluster/pkg/random"

	"net"
	"sync"
)

type Strategy string

const (
	CREATING    int64 = iota
	STARTING          = iota
	RUNNING           = iota
	TERMINATING       = iota
	TERMINATED        = iota
	FAILED            = iota

	RANDOM     Strategy = "RANDOM"
	ROUNDROBIN Strategy = "ROUNDROBIN"
)

type Proxy struct {
	Strategy Strategy
	Source   string
	Targets  []string

	listener       net.Listener
	lock           sync.Mutex
	conns          []net.Conn
	wg             sync.WaitGroup
	state          int64
	targetConnMap  *TargetConnMap
	lastUsedTarget string
}

func New() *Proxy {
	var p Proxy
	p.targetConnMap = newTargetConnMap()
	p.state = CREATING
	return &p
}

func (p *Proxy) nextTarget() string {
	p.lock.Lock()
	defer func() { p.lock.Unlock() }()
	var target string
	switch p.Strategy {
	case RANDOM:
		target = p.Targets[int(random.Get()*float32(len(p.Targets)))]
	case ROUNDROBIN:
		fallthrough
	default:
		hitLastUsed := false
		for _, t := range p.Targets {
			if hitLastUsed {
				target = t
				p.lastUsedTarget = t
				break
			} else if t == p.lastUsedTarget {
				hitLastUsed = true
			}
		}
		if target == "" {
			target = p.Targets[0]
		}
		if target == "" {
			panic("failed to select a target for proxy")
		}
	}
	return target
}
