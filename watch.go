package main

import (
	"fmt"
	"regexp"
	"time"

	PollWatch "github.com/radovskyb/watcher"
)

// Watch holds the variables to watch files
type Watch struct {
	watcher *PollWatch.Watcher
}

// NewWatch returns a new Watch struct
func NewWatch() *Watch {
	w := &Watch{
		watcher: PollWatch.New(),
	}

	w.watcher.SetMaxEvents(1)

	w.watcher.FilterOps(PollWatch.Write)

	return w
}

// SetDirRecursive sets the recurisive directory watch
func (w *Watch) SetDirRecursive(dir string) error {
	err := w.watcher.AddRecursive(dir)
	if err != nil {
		return err
	}
	return nil
}

// SetRegexFilter used to filter out specific files
func (w *Watch) SetRegexFilter(reg *regexp.Regexp) {
	w.watcher.AddFilterHook(PollWatch.RegexFilterHook(reg, false))
}

// Watch starts the watch process
func (w *Watch) Watch() error {
	go func() {
		for {
			select {
			case event := <-w.watcher.Event:
				if event.Op == PollWatch.Write {
					if trace == true {
						fmt.Println(event)
					}
					DelveEvent("restart")
				}
			case err := <-w.watcher.Error:
				if verbose == true {
					fmt.Println("[DelveWatch] watcher error: \"", err.Error(), "\"")
				}
				return
			case <-w.watcher.Closed:
				return
			}
		}
	}()

	// Start the watching process - it'll check for changes every 1 second
	if err := w.watcher.Start(1 * time.Second); err != nil {
		return err
	}

	return nil
}
