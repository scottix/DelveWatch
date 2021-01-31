package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Fprintln(os.Stdout, "Echo 1 Stdout")
	fmt.Fprintln(os.Stderr, "Echo 2 Stderr")
	time.Sleep(20 * time.Second)
	fmt.Println("Breakpoint here")
}
