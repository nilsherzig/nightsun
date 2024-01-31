package main

type YamlConfig struct {
	Modules []struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Producer    string   `yaml:"producer,omitempty"`
		Consumer    []string `yaml:"consumer,omitempty"`
		Alias       string   `yaml:"alias,omitempty"`
	} `yaml:"modules"`
}
