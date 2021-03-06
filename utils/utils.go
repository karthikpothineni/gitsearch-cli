package utils

import (
	"errors"
	"reflect"
	"strings"
)

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

// ValidateOptions - validates the options
func ValidateOptions(orgName, authKey string) error {
	if strings.TrimSpace(orgName) == "" {
		return errors.New("error: organization name cannot be empty")
	}

	if strings.TrimSpace(authKey) == "" {
		return errors.New("error: auth key cannot be empty")
	}
	return nil
}

// RemoveDuplicates - removes duplicate elements from string slice
func RemoveDuplicates(input []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}
