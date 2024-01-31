package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pkg/errors"
)

func main() {
	config, err := parseConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	items := []Item{}
	itemC := make(chan Item)
	reports := []Report{}
	reportC := make(chan Report)

	defer func() {

	}()

	// collector
	go func() {
		for {
			select {
			case item, ok := <-itemC:
				if !ok {
					return
				}

				items = append(items, item)
			case report, ok := <-reportC:
				if !ok {
					continue
				}

				reports = append(reports, report)
			}
		}
	}()

	// stop collector when all producers are done
	go func() {
		wg.Wait()
		close(itemC)
	}()

	for _, module := range config.Modules {
		if len(module.Producer) <= 0 {
			continue
		}

		wg.Add(1)
		go func(module *Module) {
			defer wg.Done()

			command := module.Producer
			cmd := exec.Command("bash", "-c", "set -euo pipefail;"+command)

			stderr, err := cmd.StderrPipe()
			if err != nil {
				reportC <- Report{
					Command: command,
					Error:   errors.Wrapf(err, "%v: could not capture stderr", module.Producer).Error(),
				}
				return
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				stderr, err := io.ReadAll(stderr)
				if err != nil {
					reportC <- Report{
						Command: command,
						Error:   errors.Wrapf(err, "%v: could not read stderr", module.Producer).Error(),
					}
					return
				}

				error := strings.TrimSpace(string(stderr))
				if len(error) <= 0 {
					return
				}

				reportC <- Report{
					Command: command,
					Error:   error,
				}
			}()

			out, err := cmd.Output()
			if err != nil {
				reportC <- Report{
					Command: command,
					Error:   errors.Wrapf(err, "%v: command failed", module.Producer).Error(),
				}
				return
			}

			trimmed := strings.TrimSpace(string(out))
			if len(trimmed) <= 0 {
				return
			}

			lines := strings.Split(trimmed, "\n")
			for _, line := range lines {
				itemC <- Item{
					Module: module,
					Line:   strings.ReplaceAll(line, "\t", "    "),
				}
			}
		}(module)
	}

	exit := 0

	if err := helper(&config, &items); err != nil {
		log.Print(err)
		exit = 1
	}

	log.Print("Producer error reports:")
	for _, report := range reports {
		log.Printf("\t%v: %v", report.Command, report.Error)
	}

	os.Exit(exit)
}

func helper(config *Config, items *[]Item) error {
	var mut sync.RWMutex

	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return (*items)[i].Module.Prefix + (*items)[i].Line
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i < 0 || i >= len(*items) {
				return ""
			}
			return (*items)[i].Module.Show()
		}),
		fuzzyfinder.WithHotReloadLock(mut.RLocker()),
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "selection failed"))
	}

	item := (*items)[idx]
	return item.Module.Exec(config, item.Line)
}
