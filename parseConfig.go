package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func parseConfigFile(path string) (Config, error) {
	yamlConfig := Config{}

	file, err := os.Open(path)
	if err != nil {
		return yamlConfig, err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&yamlConfig)
	if err != nil {
		return yamlConfig, err
	}

	return yamlConfig, nil
}

func parseConfig(filename string) (Config, error) {
	possibleConfigPaths := []string{
		filepath.Join(os.Getenv("HOME"), ".config", "unified-search", filename),
		filepath.Join(os.Getenv("HOME"), ".config", filename),
		filepath.Join(os.Getenv("HOME"), filename),
		filename,
	}

	for _, path := range possibleConfigPaths {
		if _, err := os.Stat(path); err == nil {
			return parseConfigFile(path)
		}
	}

	return Config{}, nil
}
