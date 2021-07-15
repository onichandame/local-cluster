package proxy

import (
	"errors"
	"sync/atomic"
)

func (p *Proxy) stop(force bool) error {
	p.lock.Lock()
	defer func() { p.lock.Unlock() }()
	if ok := atomic.CompareAndSwapInt64(&p.state, RUNNING, TERMINATING); !ok {
		return errors.New("cannot terminate a non-running proxy!")
	}
	var err error
	defer func() {
		var finalState int64
		if err == nil {
			finalState = TERMINATED
		} else {
			finalState = FAILED
		}
		if ok := atomic.CompareAndSwapInt64(&p.state, TERMINATING, finalState); !ok {
			panic("proxy state changed before Stop action finished! must be a race condition!")
		}
	}()
	p.wg.Done()
	if force {
		// if force, terminate all connections brutally
		if err = p.targetConnMap.terminate(); err != nil {
			return err
		}
	} else {
		// if not force, wait for all connection to finish
		p.wg.Wait()
	}
	if err = p.listener.Close(); err != nil {
		return err
	}
	return err
}

func (p *Proxy) Terminate() error {
	return p.stop(true)
}

func (p *Proxy) Shutdown() error {
	return p.stop(false)
}
