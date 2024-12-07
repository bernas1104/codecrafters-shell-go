package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BuiltIn func(string, []string)

var builtIns = make(map[string]BuiltIn)

const COMMAND_NOT_FOUND = "command not found"

func main() {
	initializeBuiltIns()

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			panic("Error!")
		}

		sanitizedInput := input[0:getInputSize(input)]
		command, args := readCommand(sanitizedInput)

		if operation, exists := builtIns[command]; exists {
			operation(command, args)
		} else {
			fmt.Printf("%v: %v\n", command, COMMAND_NOT_FOUND)
		}
	}
}

func initializeBuiltIns() {
	builtIns["exit"] = exit
	builtIns["echo"] = echo
	builtIns["type"] = typeFunc
}

func exit(_ string, args []string) {
	code, err := strconv.Atoi(args[0])

	if err != nil {
		panic("Invalid argument. Exit code must be an integer")
	}

	os.Exit(code)
}

func echo(_ string, args []string) {
	for idx, arg := range args {
		if idx == len(args)-1 {
			fmt.Printf("%v\n", arg)
		} else {
			fmt.Printf("%v ", arg)
		}
	}
}

func typeFunc(command string, args []string) {
	if _, exists := builtIns[args[0]]; exists {
		fmt.Printf("%v is a shell builtin\n", args[0])
	} else {
		fmt.Printf("%v: not found\n", args[0])
	}
}

func getInputSize(input string) int {
	return len(input) - 1
}

func readCommand(command string) (string, []string) {
	var splitedCommand = strings.Split(command, " ")

	return splitedCommand[0], splitedCommand[1:]
}
