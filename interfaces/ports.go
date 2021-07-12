package interfaces

import "sync"

var LockedPortsMap = make(map[uint]interface{})

var LastRegisteredPort uint

var Lock sync.Mutex
