package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type BuiltIn func([]string)

var builtIns = make(map[string]BuiltIn)
var paths = strings.Split(os.Getenv("PATH"), ":")

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
			operation(args)
			continue
		}

		commandExecuted := false
		for _, path := range paths {
			if executableExistsInPath(command, path) {
				cmd := exec.Command(getExecutablePath(command, path), args...)
				commandExecuted = true

				output, err := cmd.Output()
				if err != nil {
					fmt.Printf("Error while executing %v", command)
					break
				}

				fmt.Printf("%v", string(output))
				break
			}
		}

		if commandExecuted {
			continue
		}

		fmt.Printf("%v: %v\n", command, COMMAND_NOT_FOUND)
	}
}

func initializeBuiltIns() {
	builtIns["exit"] = exit
	builtIns["echo"] = echo
	builtIns["type"] = typeFunc
}

func getInputSize(input string) int {
	return len(input) - 1
}

func readCommand(command string) (string, []string) {
	var splitedCommand = strings.Split(command, " ")

	return splitedCommand[0], splitedCommand[1:]
}

func exit(args []string) {
	code, err := strconv.Atoi(args[0])

	if err != nil {
		panic("Invalid argument. Exit code must be an integer")
	}

	os.Exit(code)
}

func echo(args []string) {
	for idx, arg := range args {
		if idx == len(args)-1 {
			fmt.Printf("%v\n", arg)
		} else {
			fmt.Printf("%v ", arg)
		}
	}
}

func typeFunc(args []string) {
	if _, exists := builtIns[args[0]]; exists {
		fmt.Printf("%v is a shell builtin\n", args[0])
		return
	}

	for _, path := range paths {
		if executableExistsInPath(args[0], path) {
			fmt.Printf("%v is %v\n", args[0], getExecutablePath(args[0], path))
			return
		}
	}

	fmt.Printf("%v: not found\n", args[0])
}

func executableExistsInPath(command string, path string) bool {
	if _, err := os.Stat(getExecutablePath(command, path)); err == nil {
		return true
	}

	return false
}

func getExecutablePath(executable string, path string) string {
	return path + "/" + executable
}
