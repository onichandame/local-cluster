package portallocator

import "testing"

func TestPortAllocator(t *testing.T) {
	a := New(10000, 20000)
	p1, err := a.Allocate()
	if err != nil {
		t.Error(err)
	}
	p2, err := a.Allocate()
	if err != nil {
		t.Error(err)
	}
	if p1 == p2 {
		t.Error("failed to allocate different ports!")
	}
}
