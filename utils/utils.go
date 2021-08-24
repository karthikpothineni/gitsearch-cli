package utils

import "reflect"

// GetMapKeys - returns keys in a map
func GetMapKeys(v interface{}) []string {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return nil
	}
	t := rv.Type()
	if t.Key().Kind() != reflect.String {
		return nil
	}
	var result []string
	for _, kv := range rv.MapKeys() {
		result = append(result, kv.String())
	}
	return result
}

// HandleNilString - return empty string if string pointer is nil
func HandleNilString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
