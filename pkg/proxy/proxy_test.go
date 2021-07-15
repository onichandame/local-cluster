package proxy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"

	portallocator "github.com/onichandame/local-cluster/pkg/port_allocator"
)

func TestProxy(t *testing.T) {
	var err error
	var src, tgt uint
	portAllocator := portallocator.New(10000, 20000)
	if src, err = portAllocator.Allocate(); err != nil {
		t.Error(err)
	}
	if tgt, err = portAllocator.Allocate(); err != nil {
		t.Error(err)
	}
	p := New()
	p.Source = fmt.Sprintf("localhost:%d", src)
	p.Targets = []string{fmt.Sprintf("localhost:%d", tgt)}
	var server net.Listener
	if server, err = net.Listen("tcp", fmt.Sprintf(":%d", tgt)); err != nil {
		t.Error(err)
	}
	defer func() {
		if server != nil {
			server.Close()
		}
	}()
	greeting := []byte("hi")
	go func() {
		for {
			if conn, err := server.Accept(); err != nil {
				t.Error(err)
			} else {
				conn.Write(greeting)
				conn.Close()
			}
		}
	}()
	if err = p.Start(); err != nil {
		t.Error(err)
	}
	var client net.Conn
	if client, err = net.Dial("tcp", fmt.Sprintf("localhost:%d", src)); err != nil {
		t.Error(err)
	}
	buf := []byte{}
	tmp := []byte{}
	for {
		n, err := client.Read(tmp)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				t.Error(err)
			}
		}
		buf = append(buf, tmp[:n]...)
	}
	if bytes.Compare(buf, greeting) != 0 {
		t.Error("failed to receive response through proxy!")
	}
}
