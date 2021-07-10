package utils

import (
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func Download(url string, path string) error {
	logrus.Infof("downloading %s to %s", url, path)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := io.Copy(file, res.Body)
	logrus.Infof("downloaded %d bytes", bytes)
	return err
}
