package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetMapKeys - test if keys are extracted when valid map is passed
func TestGetMapKeys(t *testing.T) {
	check := assert.New(t)
	inputMap := make(map[string]string)
	inputMap["google"] = "golang"

	keys := GetMapKeys(inputMap)
	check.EqualValues([]string{"google"}, keys)
}

// TestGetMapKeysEmptyMap - test if keys are extracted when empty map is passed
func TestGetMapKeysEmptyMap(t *testing.T) {
	check := assert.New(t)
	inputMap := make(map[string]string)

	keys := GetMapKeys(inputMap)
	check.EqualValues([]string(nil), keys)
}

// TestGetMapKeysNilMap - test if keys are extracted when nil is passed
func TestGetMapKeysNilMap(t *testing.T) {
	check := assert.New(t)

	keys := GetMapKeys(nil)
	check.EqualValues([]string(nil), keys)
}

// TestGetMapKeysIntKey - test if keys are extracted when key is integer
func TestGetMapKeysIntKey(t *testing.T) {
	check := assert.New(t)
	inputMap := make(map[int]string)
	inputMap[123] = "golang"

	keys := GetMapKeys(inputMap)
	check.EqualValues([]string(nil), keys)
}

// TestHandleNilString - test if valid string is returned when string pointer is passed
func TestHandleNilString(t *testing.T) {
	check := assert.New(t)

	input := "go-github"
	result := HandleNilString(&input)
	check.Equal(input, result)
}

// TestHandleNilStringNilPointer - test if valid string is returned when string pointer is nil
func TestHandleNilStringNilPointer(t *testing.T) {
	check := assert.New(t)

	result := HandleNilString(nil)
	check.Equal("", result)
}

// TestValidateOptionsNoError - test options when there is no error
func TestValidateOptionsNoError(t *testing.T) {
	check := assert.New(t)
	orgName := "golang"
	authKey := "ghp_jjkhwMWqXJzWethZdCsyuZjCAYaV072BznRD"

	err := ValidateOptions(orgName, authKey)
	check.Nil(err)
}

// TestValidateOptionsOrgNameError - test options when organization name is invalid
func TestValidateOptionsOrgNameError(t *testing.T) {
	check := assert.New(t)
	authKey := "ghp_jjkhwMWqXJzWethZdCsyuZjCAYaV072BznRD"

	err := ValidateOptions("", authKey)
	check.Contains(err.Error(), "organization name cannot be empty")
}

// TestValidateOptionsAuthKeyError - test options when auth key is invalid
func TestValidateOptionsAuthKeyError(t *testing.T) {
	check := assert.New(t)
	orgName := "golang"

	err := ValidateOptions(orgName, "")
	check.Contains(err.Error(), "auth key cannot be empty")
}
