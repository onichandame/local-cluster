package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type HashAlgorithm string

const (
	Sha256 HashAlgorithm = "sha256"
)

func CheckFileHash(filePath string, hashStr string) error {
	algorithm, hash, err := parseHashString(hashStr)
	if err != nil {
		return err
	}
	if !isAlgorithmSupported(algorithm) {
		return errors.New(fmt.Sprintf("hash algorithm %s not supported", algorithm))
	}
	switch algorithm {
	case Sha256:
		if checksum, err := checksumFileSha256(filePath); err != nil {
			return err
		} else {
			if checksum != hash {
				return errors.New("checksum not correct")
			} else {
				return nil
			}
		}
	default:
		return errors.New(fmt.Sprintf("hash algorithm %s not supported", algorithm))
	}
}

func isAlgorithmSupported(algorithm HashAlgorithm) bool {
	validator := map[HashAlgorithm]interface{}{Sha256: nil}
	_, ok := validator[algorithm]
	return ok
}

func parseHashString(hashStr string) (HashAlgorithm, string, error) {
	arr := strings.Split(hashStr, ":")
	if len(arr) < 1 {
		return "", "", errors.New("hash string cannot be nil")
	} else if len(arr) == 1 {
		return "", arr[0], nil
	} else if len(arr) > 2 {
		return "", "", errors.New("hash string broken!")
	} else {
		return HashAlgorithm(arr[0]), arr[1], nil
	}
}

func checksumFileSha256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return string(hasher.Sum(nil)), nil
}
