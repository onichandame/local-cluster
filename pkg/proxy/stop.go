package proxy

import (
	"errors"
	"sync/atomic"
)

func (p *Proxy) stop(force bool) error {
	p.lock.Lock()
	defer func() { p.lock.Unlock() }()
	if ok := atomic.CompareAndSwapInt64(&p.state, RUNNING, TERMINATING); !ok {
		return errors.New("cannot terminate a not-running proxy!")
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
	if !force {
	}
	return err
}
