package common_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Compare an expected json string with the json value of an interface. Ignores whitespaces.
func CompareJsonObjectWithJsonString(t *testing.T, expected string, actual interface{}) bool {
	jsonBytes, err := json.Marshal(actual)
	assert.Nil(t, err)
	jsonText := string(jsonBytes)
	return jsonText == remarshal(expected)
}

// Unmarshal object and Marshal it again to remove white spaces and new lines
func remarshal(input string) string {
	var dynamicObj interface{}
	err := json.Unmarshal([]byte(input), &dynamicObj)

	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.Marshal(dynamicObj)
	if err != nil {
		panic(err)
	}

	return string(jsonBytes)
}
