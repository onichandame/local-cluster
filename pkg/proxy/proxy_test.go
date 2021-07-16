package proxy

import (
	"errors"
	"fmt"
	"io"
	"net"
	"testing"

	portallocator "github.com/onichandame/local-cluster/pkg/port_allocator"
)

func TestProxy(t *testing.T) {
	var err error
	var src, tgt1, tgt2 uint
	portAllocator := portallocator.New(10000, 20000)
	if src, err = portAllocator.Allocate(); err != nil {
		t.Error(err)
	}
	if tgt1, err = portAllocator.Allocate(); err != nil {
		t.Error(err)
	}
	if tgt2, err = portAllocator.Allocate(); err != nil {
		t.Error(err)
	}
	p := New()
	p.Strategy = ROUNDROBIN
	p.Source = fmt.Sprintf("localhost:%d", src)
	p.Targets = []string{fmt.Sprintf("localhost:%d", tgt1), fmt.Sprintf("localhost:%d", tgt2)}
	if err = p.Start(); err != nil {
		t.Error(err)
	}
	defer func() { p.Terminate() }()
	var server1, server2 net.Listener
	if server1, err = net.Listen("tcp", fmt.Sprintf(":%d", tgt1)); err != nil {
		t.Error(err)
	}
	if server2, err = net.Listen("tcp", fmt.Sprintf(":%d", tgt2)); err != nil {
		t.Error(err)
	}
	defer func() {
		if server1 != nil {
			server1.Close()
		}
		if server2 != nil {
			server2.Close()
		}
	}()
	greeting1 := []byte("hi1")
	greeting2 := []byte("hi2")
	go func() {
		for {
			if conn, err := server1.Accept(); err != nil {
				break
			} else {
				conn.Write(greeting1)
				conn.Close()
			}
		}
	}()
	go func() {
		for {
			if conn, err := server2.Accept(); err != nil {
				break
			} else {
				conn.Write(greeting2)
				conn.Close()
			}
		}
	}()
	readAll := func(conn net.Conn) []byte {
		buf := make([]byte, 4096)
		tmp := make([]byte, 4096)
		for {
			n, err := conn.Read(tmp)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else {
					t.Error(err)
				}
			}
			buf = append(buf, tmp[:n]...)
		}
		return buf
	}
	request := func(addr string) string {
		var client net.Conn
		if client, err = net.Dial("tcp", fmt.Sprintf("localhost:%d", src)); err != nil {
			t.Error(err)
		} else {
			defer func() { client.Close() }()
		}
		return string(readAll(client))
	}
	if string(request(fmt.Sprintf("localhost:%d", src))) == string(greeting1) {
		t.Error("failed to receive response through proxy!")
	}
	if string(request(fmt.Sprintf("localhost:%d", src))) == string(greeting2) {
		t.Error("failed to receive response through proxy!")
	}
}
