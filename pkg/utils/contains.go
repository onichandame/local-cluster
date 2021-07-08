package utils

func StrSliceToIfSlice(raw []string) []interface{} {
	res := make([]interface{}, len(raw))
	for _, s := range raw {
		res = append(res, s)
	}
	return res
}

func Contains(slice []interface{}, item interface{}) bool {
	mapper := make(map[interface{}]interface{}, len(slice))
	for _, i := range slice {
		mapper[i] = nil
	}
	_, ok := mapper[item]
	return ok
}
