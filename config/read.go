package config

import (
	"errors"
	"fmt"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/utils"
	"gopkg.in/yaml.v2"
	"os"
)

// Read YAML Configuration
func Read(path string) (Configuration, error) {
	yamlConfig, err := utils.ReadFileToString(path)
	if os.IsNotExist(err) {
		return Configuration{}, errors.New(fmt.Sprintf("the configuration file '%s' does not exist", path))
	}

	// Parse Configuration
	config := Configuration{}

	err = yaml.Unmarshal([]byte(yamlConfig), &config)
	if err != nil {
		return Configuration{}, err
	}

	return config, nil
}
