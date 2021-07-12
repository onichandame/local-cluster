package utils

import (
	"fmt"
	"net"
)

func IsPortAvailable(p uint) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
	if err == nil {
		listener.Close()
		return true
	} else {
		return false
	}
}
