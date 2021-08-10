package utils

import "github.com/sirupsen/logrus"

func RecoverAndLog() {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	defer RecoverFromError(&err)
}
