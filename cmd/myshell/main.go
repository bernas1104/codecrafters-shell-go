package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const COMMAND_NOT_FOUND = "command not found"

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		read, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			panic("Error!")
		}

		command := read[0:getReadSize(read)]

		if command == "exit 0" {
			break
		}

		fmt.Printf("%v: %v\n", command, COMMAND_NOT_FOUND)
	}
}

func getReadSize(read string) int {
	return len(read) - 1
}
