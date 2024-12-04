package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	FromEnv                  string            `key:"ENV_ARG"`
	FromSettingsFile         string            `yaml:"fromSettings"`
	FromCommandLineArguments string            `key:"CLA"`
	Sub                      *SubConfiguration `yaml:"sub" key:"SUB"`
}

type SubConfiguration struct {
	FromEnv                  string `key:"ENV_ARG"`
	FromCommandLineArguments string `key:"CLI_ARGUMENT"`
	FromSettingsFile         string `yaml:"fromSettings"`
}

func Load(configFile string) (*Configuration, error) {
	cfg := &Configuration{}

	loadDefaults(cfg)
	err := loadFromSettingsFile(cfg, configFile)
	if err != nil {
		return nil, err
	}
	loadFromEnvironmentVariables(cfg)
	loadFromCommandLineArguments(cfg)

	return cfg, nil
}

func loadDefaults(cfg *Configuration) {

}

func loadFromSettingsFile(cfg *Configuration, configFile string) error {
	ext := filepath.Ext(configFile)

	switch {
	case ext == ".yaml" || ext == ".yml":
		yamlFile, err := os.Open(configFile)
		if err != nil {
			return err
		}
		defer yamlFile.Close()

		bytes, err := io.ReadAll(yamlFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(bytes, cfg)
		if err != nil {
			return err
		}
	case ext == ".json":
		jsonFile, err := os.Open(configFile)
		if err != nil {
			return err
		}
		defer jsonFile.Close()

		bytes, err := io.ReadAll(jsonFile)
		if err != nil {
			return err
		}

		err = json.Unmarshal(bytes, cfg)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported settings file extension: %s", ext)
	}

	return nil
}

func loadFromEnvironmentVariables(cfg *Configuration) {

}

func loadFromCommandLineArguments(cfg *Configuration) {

}
