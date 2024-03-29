package main

import (
	"bytes"
	"fmt"
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
					Error:   errors.Wrap(err, "could not capture stderr").Error(),
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
						Error:   errors.Wrap(err, "could not read stderr").Error(),
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
					Error:   errors.Wrap(err, "command failed").Error(),
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
					Line:   line,
					Show:   strings.ReplaceAll(line, "\t", "    "),
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
			if (*items)[i].Module.Prefix == "" {
				return (*items)[i].Show
			}
			return fmt.Sprintf("(%s) ", (*items)[i].Module.Prefix) + (*items)[i].Show
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i < 0 || i >= len(*items) {
				return ""
			}
			return (*items)[i].Module.Show()
		}),
		fuzzyfinder.WithHotReloadLock(mut.RLocker()),
		fuzzyfinder.WithHeader(">> May the Nightsun illuminate your path <<"),
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "selection failed"))
	}

	item := (*items)[idx]
	script := strings.Join([]string{
		"set -xeuo pipefail",
		config.Modules.MkScript(),
		item.Module.Name,
	}, "\n")

	cmd := exec.Command("bash")
	cmd.Env = append(os.Environ(), "sel="+item.Line)
	cmd.Stdin = bytes.NewReader([]byte(script))
	log.Print(script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "could not run command")
	}

	return nil
}
