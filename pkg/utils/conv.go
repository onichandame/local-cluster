package utils

import "strconv"

func StrSliceToIfSlice(raw []string) []interface{} {
	res := make([]interface{}, len(raw))
	for _, s := range raw {
		res = append(res, s)
	}
	return res
}

func UintSliceToIfSlice(raw []uint) []interface{} {
	res := make([]interface{}, len(raw))
	for _, i := range raw {
		res = append(res, i)
	}
	return res
}

func StrToUint(raw string) (u uint, err error) {
	defer RecoverFromError(&err)
	if i, err := strconv.Atoi(raw); err == nil {
		u = uint(i)
	} else {
		panic(err)
	}
	return u, err
}
