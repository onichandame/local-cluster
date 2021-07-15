package proxy

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

type Config struct {
	Strategy Strategy
}

func (p *Proxy) Start() error {
	var err error
	p.lock.Lock()
	defer func() { p.lock.Unlock() }()
	defer func() {
		var finalState int64
		if err == nil {
			finalState = RUNNING
		} else {
			finalState = FAILED
		}
		if ok := atomic.CompareAndSwapInt64(&p.state, STARTING, finalState); !ok {
			panic(fmt.Sprintf("proxy state changed before Start action finished! must have been a race condition!"))
		}
	}()
	if ok := atomic.CompareAndSwapInt64(&p.state, CREATING, STARTING); !ok {
		return errors.New("proxy already started or terminated cannot be started again!")
	}
	if err = parseTCPAddr(&p.Source); err != nil {
		return err
	}
	if p.listener, err = net.Listen("tcp", p.Source); err != nil {
		return err
	}
	p.wg.Add(1)
	handleRequest := func(source net.Conn) {
		p.wg.Add(1)
		defer func() { p.wg.Done() }()
		defer func() { source.Close() }()
		targetAddr := p.nextTarget()
		if err = parseTCPAddr(&targetAddr); err != nil {
			panic("target address " + targetAddr + " not valid")
		}
		var target net.Conn
		if target, err = net.Dial("tcp", targetAddr); err != nil {
			panic(err)
		} else {
			go io.Copy(source, target)
			io.Copy(target, source)
		}
	}
	go func() {
		for {
			if conn, err := p.listener.Accept(); err != nil {
				panic(err)
			} else {
				if p.state != RUNNING {
					conn.Write([]byte("cannot handle new requests when proxy is not ready"))
					conn.Close()
				}
				go handleRequest(conn)
			}
		}
	}()
	return err
}
