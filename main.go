package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	config, err := parseConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	outputs := Outputs{}
	for _, module := range config.Modules {
		if len(module.Producer) > 0 {
			cmd := exec.Command("bash", "-c", module.Producer)
			out, err := cmd.Output()
			if err != nil {
				log.Print(errors.Wrapf(err, "%v: command failed", module.Producer))
				continue
			}

			outputs = append(outputs, Output{
				Module: module,
				Lines:  strings.Split(string(out), "\n"),
			})
		}
	}

	log.Printf(
		"got %v total lines from %v successful commands",
		len(outputs.Lines()), len(outputs))

	// items := []list.Item{}
	// for _, module := range yamlConfig.Modules {
	// 	// if len(module.Consumer) > 0 {
	// 	// 	items = append(items, item{title: module.Name, desc: module.Consumer[0]})
	// 	// 	continue
	// 	// }
	// 	// items = append(items, item{title: module.Name, desc: module.Description})
	// 	items = append(items, item{
	// 		title: module.Name,
	// 		desc:  module.Desc,
	// 	})
	// }

	// m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	// m.list.Title = "unified-search"

	// p := tea.NewProgram(m, tea.WithAltScreen())

	// if _, err := p.Run(); err != nil {
	// 	fmt.Println("Error running program:", err)
	// 	os.Exit(1)
	// }
}
