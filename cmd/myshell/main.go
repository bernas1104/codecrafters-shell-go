package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const COMMAND_NOT_FOUND = "command not found"

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			panic("Error!")
		}

		sanitizedInput := input[0:getInputSize(input)]
		command, args := readCommand(sanitizedInput)

		if command == "exit" {
			code, err := strconv.Atoi(args[0])

			if err != nil {
				panic("Invalid argument. Exit code must be an integer")
			}

			os.Exit(code)
		} else if command == "echo" {
			for idx, arg := range args {
				if idx == len(args)-1 {
					fmt.Printf("%v", arg)
				} else {
					fmt.Printf("%v ", arg)
				}
			}

			fmt.Println()
		} else {
			fmt.Printf("%v: %v\n", command, COMMAND_NOT_FOUND)
		}
	}
}

func getInputSize(input string) int {
	return len(input) - 1
}

func readCommand(command string) (string, []string) {
	var splitedCommand = strings.Split(command, " ")

	return splitedCommand[0], splitedCommand[1:]
}
