package proxy

import (
	"net"
	"sync"
)

type Conns struct {
	lock  sync.Mutex
	conns map[net.Conn]interface{}
}

func newConns() Conns {
	var m Conns
	m.conns = make(map[net.Conn]interface{})
	return m
}

func (m *Conns) addConn(conn net.Conn) {
	m.lock.Lock()
	defer func() { m.lock.Unlock() }()
	m.conns[conn] = nil
}

func (m *Conns) delConn(conn net.Conn) {
	m.lock.Lock()
	defer func() { m.lock.Unlock() }()
	delete(m.conns, conn)
}

func (m *Conns) listConns() []net.Conn {
	m.lock.Lock()
	defer func() { m.lock.Unlock() }()
	var conns []net.Conn
	for conn := range m.conns {
		conns = append(conns, conn)
	}
	return conns
}

type TargetConnMap struct {
	lock    sync.Mutex
	targets map[string]Conns
}

func newTargetConnMap() TargetConnMap {
	var m TargetConnMap
	m.targets = make(map[string]Conns)
	return m
}

func (m *TargetConnMap) initTargets(targets []string) {
	m.lock.Lock()
	defer func() { m.lock.Unlock() }()
	for _, key := range targets {
		m.targets[key] = newConns()
	}
}

func (m *TargetConnMap) terminate() error {
	var err error
	for _, target := range m.targets {
		for conn := range target.conns {
			if e := conn.Close(); e != nil {
				if err == nil {
					err = e
				}
			}
		}
	}
	return err
}
