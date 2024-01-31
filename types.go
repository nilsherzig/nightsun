package main

import (
	"fmt"
)

type Module struct {
	Name     string   `yaml:"name"`
	Desc     string   `yaml:"desc"`
	Prefix   string   `yaml:"prefix"`
	Producer string   `yaml:"producer"`
	Consumer []string `yaml:"consumer"`
	Alias    string   `yaml:"alias"`
}

func (m *Module) Show() string {
	return fmt.Sprintf("%v\n%v", m.Name, m.Desc)
}

type Config struct {
	Modules []*Module `yaml:"modules"`
}

////////////////////////////////////////////////////////////////////////////////

type Item struct {
	Module *Module
	Line   string
}
