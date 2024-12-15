package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/lexer"
	"github.com/codecrafters-io/shell-starter-go/token"
)

type BuiltIn func([]string)

var builtIns = make(map[string]BuiltIn)
var paths = strings.Split(os.Getenv("PATH"), ":")

var cdPath string = ""

const COMMAND_NOT_FOUND = "command not found"

func main() {
	initializeBuiltIns()
	token.InitializeEscapeCharacters()
	currentDirectory, err := os.Getwd()

	if err != nil {
		panic("Cannot determine current directory paath")
	}

	cdPath = currentDirectory

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			panic("Error!")
		}

		sanitizedInput := input[0:getInputSize(input)]
		tokens := getTokens(sanitizedInput)

		command := token.GetCommand(tokens)

		var args []string
		if len(tokens) > 2 {
			args = token.GetArguments(tokens[2:])
		}

		if operation, exists := builtIns[command]; exists {
			operation(args)
			continue
		}

		commandExecuted := tryExecuteCommand(command, args)
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
	builtIns["pwd"] = pwd
	builtIns["cd"] = cd
}

func getInputSize(input string) int {
	return len(input) - 1
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

func tryExecuteCommand(command string, args []string) bool {
	commandExecuted := false
	for _, path := range paths {
		if executableExistsInPath(command, path) {
			cmd := exec.Command(getExecutablePath(command, path), args...)
			commandExecuted = true

			output, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
				fmt.Printf("Error while executing %v\n", command)
				break
			}

			fmt.Printf("%v", string(output))
			break
		}
	}

	return commandExecuted
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

func pwd(_ []string) {
	fmt.Println(cdPath)
}

func cd(args []string) {
	tokens := getTokens(args[0])
	firstToken := tokens[0]

	pathString := getPathString(tokens)

	if firstToken.Type == token.GO_BACK_PATH || firstToken.Type == token.RELATIVE_PATH {
		pathString = cdPath + "/" + pathString
	}

	if firstToken.Type == token.USER_PATH {
		home := os.Getenv("HOME")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}

		pathString = home + pathString[1:]
	}

	info, err := os.Stat(pathString)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("cd: %v: No such file or directory\n", args[0])
		} else {
			fmt.Println("Error: ", err)
		}

		return
	}

	absPath, err := filepath.Abs(pathString)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if info.IsDir() {
		cdPath = absPath
		return
	}

	fmt.Printf("cd: %v: is not a directory\n", absPath)
}

func getTokens(arg string) []token.Token {
	l := lexer.New(arg)
	var tokens []token.Token

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		tokens = append(tokens, tok)
	}

	return tokens
}

func getPathString(tokens []token.Token) string {
	var pathString string = ""
	for _, tok := range tokens {
		pathString += tok.Literal
	}

	return pathString
}
