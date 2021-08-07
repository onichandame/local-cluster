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
	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()
	gzReader, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}
	defer gzReader.Close()
	tarReader := tar.NewReader(gzReader)
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
		switch entry.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(entryPath, os.ModeDir); err != nil {
				return err
			}
		case tar.TypeReg:
			entryFile, err := createEntryFile(entryPath)
			if err != nil {
				return err
			}
			defer entryFile.Close()
			io.Copy(entryFile, tarReader)
		default:
			return errors.New(fmt.Sprintf("extract unknown type %d", entry.Typeflag))
		}
	}
	return nil
}

func createEntryFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModeDir); err != nil {
		return nil, err
	}
	if PathExists(path) {
		if err := os.Remove(path); err != nil {
			panic(err)
		}
	}
	return os.Create(path)
}
