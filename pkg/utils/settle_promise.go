package utils

import "github.com/chebyrash/promise"

func SettlePromise(resolve func(promise.Any), reject func(error)) {
	var err error
	RecoverFromError(&err)
	if err == nil {
		resolve(nil)
	} else {
		reject(err)
	}
}
