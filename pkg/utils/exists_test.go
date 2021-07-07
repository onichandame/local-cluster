package utils

import (
	"testing"
)

func TestFileExists(t *testing.T) {
	const fileName = "exists_test.go"
	const noFileName = "nofile"
	if !PathExists(fileName) {
		t.Error("cannot find an existing file")
	}
	if PathExists(noFileName) {
		t.Error("finds a non-existing file")
	}
}
