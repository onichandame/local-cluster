package portallocator

import (
	"errors"
	"fmt"
	"net"
	"testing"
)

func TestCheckPort(t *testing.T) {
	var start, end uint = 10000, 20000
	maxTrials := end - start
	var freePort uint
	var trial int
	for i := start; i < end; i++ {
		if trial >= int(maxTrials) {
			t.Error("cannot execute test as no free ports available!")
		}
		if err := checkPort(i); err == nil {
			freePort = i
		}
		trial++
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", freePort))
	defer func() {
		if listener != nil {
			listener.Close()
		}
	}()
	if err != nil {
		t.Error(err)
	}
	if err := checkPort(freePort); err == nil {
		t.Error(errors.New("cannot recognize port in use"))
	}
}
