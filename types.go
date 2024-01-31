package main

type Module struct {
	Name     string   `yaml:"name"`
	Desc     string   `yaml:"desc"`
	Producer string   `yaml:"producer"`
	Consumer []string `yaml:"consumer"`
	Alias    string   `yaml:"alias"`
}

type Config struct {
	Modules []Module `yaml:"modules"`
}

type Output struct {
	Module Module
	Lines  []string
}

type Outputs []Output

func (o Outputs) Lines() []string {
	lines := []string{}
	for _, output := range o {
		lines = append(lines, output.Lines...)
	}
	return lines
}
