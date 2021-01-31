package main

import (
	"flag"
	"fmt"
	"regexp"
	"sync"
)

var watch *Watch
var dlvFile string
var dlvListen string
var dlvAPI string
var dlvArgs string
var dlvStopTime int
var verbose bool
var trace bool

func main() {
	flag.StringVar(&dlvFile, "delve", "", "Specify the delve command to run")
	flag.StringVar(&dlvListen, "listen", ":2345", "Specify the delve listen variable --listen=<input>")
	flag.StringVar(&dlvAPI, "api", "2", "Specify the delve api version --api-version=<input>")
	flag.StringVar(&dlvArgs, "args", "", "Additional args to program")
	flag.IntVar(&dlvStopTime, "stoptime", 0, "Number of seconds to wait after stopping")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output when running")
	flag.BoolVar(&trace, "trace", false, "Trace debug")
	flag.Parse()

	// Watches os signals
	go Signal()

	var wg sync.WaitGroup
	wg.Add(1)
	go DelveRun(&wg)

	watch = NewWatch()
	watch.SetDirRecursive(".")
	goFiles := regexp.MustCompile(".go$")
	watch.SetRegexFilter(goFiles)

	// Starts the file watcher
	err := watch.Watch()
	if err != nil {
		fmt.Println("[DelveWatch] Watch error ", err.Error())
	}

	// Waiting for delve to finish
	wg.Wait()

	if verbose == true {
		fmt.Println("[DelveWatch] Exited cleanly")
	}
}
