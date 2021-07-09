package utils

func Contains(slice []interface{}, item interface{}) bool {
	mapper := make(map[interface{}]interface{}, len(slice))
	for _, i := range slice {
		mapper[i] = nil
	}
	_, ok := mapper[item]
	return ok
}
