package utils

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func DecompressTarGZ(archivePath string, dest string) error {
	if !PathExists(dest) {
	}
	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()
	archReader, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}
	defer archReader.Close()
	tarReader := tar.NewReader(archReader)
	for {
		entry, err := tarReader.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return err
			}
		}
		entryPath := filepath.Join(dest, entry.Name)
		if entryPath == "" {
			return errors.New(fmt.Sprintf("failed to get path for entry %s", entry.Name))
		}
		entryFile, err := createEntryFile(entryPath)
		if err != nil {
			return err
		}
		defer entryFile.Close()
		io.Copy(entryFile, tarReader)
	}
	return nil
}

func createEntryFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModeDir); err != nil {
		return nil, err
	}
	return os.Create(path)
}
