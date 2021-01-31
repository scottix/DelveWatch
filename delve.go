package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var delveCmd *exec.Cmd
var listener chan string
var stdout io.ReadCloser
var stderr io.ReadCloser

// DelveStart sets up and start the delve command
func DelveStart() error {
	if trace == true {
		fmt.Println("[DelveWatch:DelveStart] Start")
	}
	var err error

	dlvParams := []string{
		"debug",
		dlvFile,
		"--listen=" + dlvListen,
		"--headless",
		"--continue",
		"--log",
		"--accept-multiclient",
		"--api-version=" + dlvAPI,
	}

	if dlvArgs != "" {
		dlvParams = append(dlvParams, "--", dlvArgs)
	}

	delveCmd = exec.Command("dlv", dlvParams...)
	if verbose == true {
		fmt.Println("[DelveWatch] Running command", delveCmd.String())
	}

	stdout, err = delveCmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err = delveCmd.StderrPipe()
	if err != nil {
		return err
	}

	err = delveCmd.Start()
	if err != nil {
		return err
	}

	if trace == true {
		fmt.Println("[DelveWatch:DelveStart] End")
	}

	return nil
}

// DelveRun main loop to watch devlve
func DelveRun(wg *sync.WaitGroup) {
	if trace == true {
		fmt.Println("[DelveWatch:DelveRun] Start")
	}
	defer wg.Done()

	pipeOut := make(chan io.ReadCloser)
	listener = make(chan string)
	command := ""

	done := false
	for done == false {
		err := DelveStart()
		if err != nil {
			fmt.Println("[DelveWatch] Failed to start program", err.Error())
		}

		go DelveOut(pipeOut)

		pipeOut <- stdout
		pipeOut <- stderr

		command = <-listener
		if command == "restart" {
			if verbose == true {
				fmt.Println("[DelveWatch] Restarting delve")
			}
		} else {
			if verbose == true {
				fmt.Println("[DelveWatch] Shutting down")
			}
			done = true
		}

		DelveStop()

		if dlvStopTime > 0 {
			if verbose == true {
				fmt.Println("[DelveWatch] Sleeping " + strconv.Itoa(dlvStopTime) + " second(s)")
			}
			time.Sleep(time.Duration(dlvStopTime) * time.Second)
		}
	}

	if trace == true {
		fmt.Println("[DelveWatch:DelveRun] End")
	}
}

// DelveStop forcefully stops the program
func DelveStop() {
	if trace == true {
		fmt.Println("[DelveWatch:DelveStop] Start")
	}
	delveCmd.Process.Kill()
	delveCmd.Wait()
	if trace == true {
		fmt.Println("[DelveWatch:DelveStop] End")
	}
}

// DelveEvent running blocker for events
func DelveEvent(msg string) {
	listener <- msg
}

// DelveOut sends the program output to terminal
func DelveOut(pipeOut <-chan io.ReadCloser) {
	if trace == true {
		fmt.Println("[DelveWatch:DelveOut] Start")
	}

	innerDone := make(chan bool)
	inner := func(pipeInner io.ReadCloser) {
		reader := bufio.NewReader(pipeInner)

		done := false
		for done == false {
			out, err := reader.ReadString('\n')
			if err != nil {
				done = true
				continue
			}
			fmt.Print(out)
		}
		innerDone <- done
	}

	outerDone := false
	for outerDone == false {
		pipe := <-pipeOut
		go inner(pipe)

		pipe = <-pipeOut
		go inner(pipe)

		outerDone = <-innerDone
	}

	if trace == true {
		fmt.Println("[DelveWatch:DelveOut] End")
	}
}
