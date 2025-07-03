package fingerprint

func mapAdd[keyType comparable, valueType any](map1 map[keyType]valueType, map2 map[keyType]valueType) map[keyType]valueType {
	var result = make(map[keyType]valueType)
	for k, v := range map1 {
		result[k] = v
	}
	for k, v := range map2 {
		result[k] = v
	}
	return result
}
