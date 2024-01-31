package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

func (i Item) FilterValue() string {
	return i.Line
}

////////////////////////////////////////////////////////////////////////////////

type ItemDelegate struct {
}

// Height implements list.ItemDelegate.
func (ItemDelegate) Height() int {
	return 1
}

// Render implements list.ItemDelegate.
func (ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	format := ""
	if index == m.Index() {
		format = "\033[1;4;32m"
	}
	fmt.Fprintf(w, "%v%v%v",
		format, item.FilterValue(), "\033[m")
}

// Spacing implements list.ItemDelegate.
func (ItemDelegate) Spacing() int {
	return 0
}

// Update implements list.ItemDelegate.
func (ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type Model struct {
	List list.Model
}

func (Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return docStyle.Render(m.List.View())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

////////////////////////////////////////////////////////////////////////////////
