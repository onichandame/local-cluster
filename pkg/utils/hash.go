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

func CheckFileHash(filePath string, hashStr string) (err error) {
	algorithm, hash, err := parseHashString(hashStr)
	if err != nil {
		panic(err)
	}
	if !isAlgorithmSupported(algorithm) {
		panic(errors.New(fmt.Sprintf("hash algorithm %s not supported", algorithm)))
	}
	switch algorithm {
	case Sha256:
		if checksum, err := checksumFileSha256(filePath); err != nil {
			panic(err)
		} else {
			if checksum != hash {
				panic(errors.New("checksum not correct"))
			}
		}
	default:
		panic(errors.New(fmt.Sprintf("hash algorithm %s not supported", algorithm)))
	}
	return err
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

func checksumFileSha256(filePath string) (hash string, err error) {
	defer RecoverFromError(&err)
	if file, err := os.Open(filePath); err != nil {
		panic(err)
	} else {
		defer file.Close()
		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			panic(err)
		}
		hash = string(hasher.Sum(nil))
	}
	return hash, nil
}
