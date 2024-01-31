package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type Consumer struct {
	IsMod string `yaml:"mod"`
	IsCmd string `yaml:"cmd"`
}

type Module struct {
	Name     string     `yaml:"name"`
	Desc     string     `yaml:"desc"`
	Prefix   string     `yaml:"prefix"`
	Producer string     `yaml:"producer"`
	Consumer []Consumer `yaml:"consumer"`
	Alias    string     `yaml:"alias"`
}

func (m *Module) Show() string {
	return fmt.Sprintf("%v\n%v", m.Name, m.Desc)
}

func (m *Module) Exec(config *Config, selection string) error {
	if len(m.Alias) > 0 {
		module := config.FindModule(m.Alias)
		if module == nil {
			return errors.Errorf("%v: invalid alias: %v", m.Name, m.Alias)
		}
		return module.Exec(config, selection)
	}

	if len(m.Consumer) <= 0 {
		return errors.Errorf("%v: no consumers", m.Name)
	}

	for i, consumer := range m.Consumer {
		if len(consumer.IsMod) > 0 {
			parts := strings.SplitN(consumer.IsMod, " ", 2)
			if len(parts) != 2 {
				return errors.Errorf(
					"%v: consumer %v: 2 parts required!",
					m.Name, i,
				)
			}

			name := parts[0]
			module := config.FindModule(name)
			if module == nil {
				return errors.Errorf(
					"%v: consumer %v: module %v not found",
					m.Name, i, name,
				)
			}

			if err := module.Exec(config, parts[1]); err != nil {
				return err
			}
		} else if len(consumer.IsCmd) > 0 {
			cmd := exec.Command("bash", "-c", consumer.IsCmd, "--", selection)
			cmd.Env = append(os.Environ(), "sel="+selection)
			log.Print(cmd.String())

			if err := cmd.Run(); err != nil {
				return errors.Wrapf(
					err, "%v: %v: command failed",
					m.Name, consumer.IsCmd,
				)
			}
		} else {
			return errors.Errorf(
				"%v: consumer %v: either mod or cmd must be set!",
				m.Name, i,
			)
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

type Config struct {
	Modules []*Module `yaml:"modules"`
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
}
