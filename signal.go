package main

import (
	"os"
	"os/signal"
	"syscall"
)

// Signal listens for program signals
func Signal() {
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGHUP)
	done := false
	for done == false {
		s := <-chanSignal
		if s == syscall.SIGINT {
			DelveEvent("stop")
			watch.watcher.Close()
			done = true
		} else {
			DelveEvent("restart")
		}
	}
}
