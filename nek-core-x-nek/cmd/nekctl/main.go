package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		fmt.Println("nekctl commands: version, help")
		return
	}
	if os.Args[1] == "version" {
		fmt.Println("NEK Core dev")
		return
	}
	fmt.Fprintln(os.Stderr, "unknown command:", os.Args[1])
	os.Exit(2)
}
