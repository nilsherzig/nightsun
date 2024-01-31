package main

import (
	"fmt"
)

type Report struct {
	Command string
	Error   string
}

////////////////////////////////////////////////////////////////////////////////

type Module struct {
	Name     string `yaml:"name"`
	Desc     string `yaml:"desc"`
	Prefix   string `yaml:"prefix"`
	Producer string `yaml:"producer"`
	Consumer string `yaml:"consumer"`
	Alias    string `yaml:"alias"`
}

func (m *Module) Show() string {
	return fmt.Sprintf("%v\n%v", m.Name, m.Desc)
}

func (m *Module) ToFunction() string {
	if len(m.Alias) > 0 {
		return fmt.Sprintf("%v(){ %v \"$@\"; }\n", m.Name, m.Alias)
	}

	return fmt.Sprintf("%v() {\n%v\n}\n", m.Name, m.Consumer)
}

////////////////////////////////////////////////////////////////////////////////

type Modules []*Module

func (m Modules) MkScript() string {
	res := ""
	for _, m := range m {
		res += m.ToFunction()
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////

type Config struct {
	Modules Modules `yaml:"modules"`
}

func (c Config) FindModule(name string) *Module {
	var res *Module
	for _, module := range c.Modules {
		if module.Name == name {
			res = module
			break
		}
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////

type Item struct {
	Module *Module
	Line   string
	Show   string
}
