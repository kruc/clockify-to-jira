package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

const (
	ErrLoaderConfigFileNotFound   = ConfigErr("Cannot load configuration from file")
	ErrLoaderInvalidConfiguration = ConfigErr("Invalid config data - Unmarshall error")
)

func LoadFromYamlFile(filePath string) (Config, error) {
	fileConfigData, err := createFileConfigSource(filePath)

	if err != nil {
		return Config{}, ErrLoaderConfigFileNotFound
	}

	config, err := Load(fileConfigData)

	if err != nil {
		return Config{}, ErrLoaderInvalidConfiguration
	}

	return config, nil
}

func Load(configDataProvider []byte) (Config, error) {

	yamlData := configDataProvider

	config := Config{}

	err := yaml.Unmarshal(yamlData, &config)

	if err != nil {
		return Config{}, ErrLoaderInvalidConfiguration
	}

	config.Workspaces = config.combineWithDefaultConfig()

	return config, nil
}

func createFileConfigSource(filePath string) ([]byte, error) {
	fileConfigData, err := os.ReadFile(filePath)

	if err != nil {
		return nil, ErrLoaderConfigFileNotFound
	}

	return fileConfigData, nil
}
