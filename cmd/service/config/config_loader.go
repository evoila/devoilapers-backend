package config

import (
	"encoding/json"
	"io/ioutil"
)

// Function will load a configuration from the provided file path and try's to map it to the RawConfiguration struct.
func LoadConfigurationFromFile(filepath string) (parsedConfig RawConfig, err error) {
	fileBytes, err := ioutil.ReadFile(filepath)

	if err != nil || json.Unmarshal(fileBytes, &parsedConfig) != nil {
		return
	}

	return
}
