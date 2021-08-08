package utils

func RecoverFromError(err *error) {
	if er := recover(); er != nil {
		if e, ok := er.(error); ok {
			*err = e
		}
	} else {
		*err = nil
	}
}
