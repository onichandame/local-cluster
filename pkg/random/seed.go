package random

import (
	"math/rand"
	"sync/atomic"
	"time"
)

var seeded uint32

func Seed() {
	if ok := atomic.CompareAndSwapUint32(&seeded, 0, 1); ok {
		rand.Seed(time.Now().UnixNano())
	}
}
