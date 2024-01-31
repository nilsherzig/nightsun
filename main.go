package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pkg/errors"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func main() {
	config, err := parseConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	items := []Item{}
	for _, module := range config.Modules {
		if len(module.Producer) > 0 {
			cmd := exec.Command("bash", "-c", "set -euo pipefail;"+module.Producer)
			out, err := cmd.Output()
			if err != nil {
				log.Print(errors.Wrapf(err, "%v: command failed", module.Producer))
				continue
			}

			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			for _, line := range lines {
				items = append(items, Item{
					Module: module,
					Line:   line,
				})
			}
		}
	}

	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return items[i].Module.Prefix + items[i].Line
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i < 0 || i >= len(items) {
				return ""
			}
			return items[i].Module.Show()
		}),
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "selection failed"))
	}

	log.Print(items[idx].Line)
	log.Print(items[idx].Module.Show())

	// m := Model{
	// 	List: list.New(outputs.Items(), ItemDelegate{}, 0, 0),
	// }

	// p := tea.NewProgram(m)
	// if _, err := p.Run(); err != nil {
	// 	log.Fatal(errors.Wrap(err, "failed to start UI"))
	// }

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
