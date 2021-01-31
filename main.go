package main

import (
	"flag"
	"fmt"
	"regexp"
)

var watch *Watch
var dlvFile string
var dlvListen string
var dlvAPI string
var dlvArgs string
var verbose bool
var trace bool

func main() {
	flag.StringVar(&dlvFile, "delve", "", "Specify the delve command to run")
	flag.StringVar(&dlvListen, "listen", ":2345", "Specify the delve listen variable --listen=<input>")
	flag.StringVar(&dlvAPI, "api", "2", "Specify the delve api version --api-version=<input>")
	flag.StringVar(&dlvArgs, "args", "", "Additional args to program")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output when running")
	flag.BoolVar(&trace, "trace", false, "Trace debug")
	flag.Parse()

	watch = NewWatch()
	watch.SetDirRecursive(".")
	goFiles := regexp.MustCompile(".go$")
	watch.SetRegexFilter(goFiles)

	go DelveRun()
	go Signal()

	err := watch.Watch()
	if err != nil {
		fmt.Println("[DelveWatch] Watch error ", err.Error())
	}

	if verbose == true {
		fmt.Println("[DelveWatch] Exited cleanly")
	}
}
