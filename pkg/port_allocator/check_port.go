package portallocator

import (
	"fmt"
	"net"
)

func checkPort(p uint) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
	if listener != nil {
		listener.Close()
	}
	return err
}
