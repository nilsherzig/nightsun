package main

import (
	"log"
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

	// collector
	go func() {
		for {
			select {
			case item, ok := <-itemC:
				if !ok {
					return
				}

				items = append(items, item)
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

			cmd := exec.Command("bash", "-c", "set -euo pipefail;"+module.Producer)
			out, err := cmd.Output()
			if err != nil {
				log.Print(errors.Wrapf(err, "%v: command failed", module.Producer))
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

	var mut sync.RWMutex

	idx, err := fuzzyfinder.Find(
		&items,
		func(i int) string {
			return items[i].Module.Prefix + items[i].Line
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i < 0 || i >= len(items) {
				return ""
			}
			return items[i].Module.Show()
		}),
		fuzzyfinder.WithHotReloadLock(mut.RLocker()),
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "selection failed"))
	}

	item := items[idx]
	if err := item.Module.Exec(&config, item.Line); err != nil {
		log.Fatal(err)
	}
}
