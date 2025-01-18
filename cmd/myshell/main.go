package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/lexer"
	"github.com/codecrafters-io/shell-starter-go/token"
	"github.com/codecrafters-io/shell-starter-go/utils"
)

type BuiltIn func([]string) []byte

var builtIns = make(map[string]BuiltIn)

var paths []string
var usrBinPath = ""
var cdPath string = ""
var redirectStdout bool = false

const COMMAND_NOT_FOUND = "command not found"

func main() {
	initializePaths()
	initializeBuiltIns()
	lexer.InitializeEscapeCharacters()
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
		tokens := lexer.GetTokens(sanitizedInput)

		var args []string
		args = lexer.GetArguments(tokens)

		command := args[0]
		args = args[1:]

		commandExecuted := tryExecuteCommand(command, args)
		if commandExecuted {
			continue
		}

		fmt.Printf("%v: %v\n", command, COMMAND_NOT_FOUND)
	}
}

func initializePaths() {
	var splitter string = ""

	if runtime.GOOS == "windows" {
		splitter = ";"
		usrBinPath = "C:\\Program Files\\Git\\usr\\bin"
	} else {
		splitter = ":"
		usrBinPath = "/usr/bin"
	}

	paths = strings.Split(os.Getenv("PATH"), splitter)
}

func initializeBuiltIns() {
	builtIns["cd"] = cd
	builtIns["exit"] = exit
	builtIns["type"] = typeFunc
	builtIns["pwd"] = pwd
	builtIns["echo"] = echo
}

func getInputSize(input string) int {
	ignoredCharacters := 2

	if utils.Contains([]string{"linux", "unix"}, runtime.GOOS) {
		ignoredCharacters = 1
	}

	return len(input) - ignoredCharacters
}

func exit(args []string) []byte {
	code, err := strconv.Atoi(args[0])

	if err != nil {
		panic("Invalid argument. Exit code must be an integer")
	}

	os.Exit(code)
	return []byte{}
}

func echo(args []string) []byte {
	cmd := exec.Command(getExecutablePath("echo", usrBinPath), args...)
	bytes, err := cmd.Output()

	if err != nil {
		fmt.Printf("%v\n", err)
		return []byte{}
	}

	return bytes
}

func typeFunc(args []string) []byte {
	if _, exists := builtIns[args[0]]; exists {
		fmt.Printf("%v is a shell builtin\n", args[0])
		return []byte{}
	}

	for _, path := range paths {
		if executableExistsInPath(args[0], path) {
			fmt.Printf("%v is %v\n", args[0], getExecutablePath(args[0], path))
			return []byte{}
		}
	}

	fmt.Printf("%v: not found\n", args[0])
	return []byte{}
}

func tryExecuteCommand(command string, args []string) bool {
	commandExecuted := false
	usedArgs := args

	if utils.AnyMultiple([]string{"1>", ">"}, args) {
		redirectStdout = true
		usedArgs = args[0:utils.FindRedirectIndex(args)]
	}

	if operation, exists := builtIns[command]; exists {
		output := operation(usedArgs)
		commandExecuted = true

		handleCommandOutput(output, args)

		return commandExecuted
	}

	for _, path := range paths {
		if executableExistsInPath(command, path) {
			cmd := exec.Command(getExecutablePath(command, path), usedArgs...)
			commandExecuted = true

			output, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
				fmt.Printf("Error while executing %v\n", command)
				break
			}

			handleCommandOutput(output, usedArgs)
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
	if runtime.GOOS == "windows" {
		return path + "\\" + executable
	}

	return path + "/" + executable
}

func handleCommandOutput(output []byte, args []string) {
	if redirectStdout {
		idx := utils.FindRedirectIndex(args)

		filePath := args[idx+1]

		err := os.WriteFile(filePath, output, 0644)
		fmt.Println(err)
		return
	}

	fmt.Printf("%v", string(output))
}

func pwd(_ []string) []byte {
	return []byte(cdPath)
}

func cd(args []string) []byte {
	tokens := lexer.GetTokens(args[0])
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

		return []byte{}
	}

	absPath, err := filepath.Abs(pathString)
	if err != nil {
		fmt.Println("Error: ", err)
		return []byte{}
	}

	if info.IsDir() {
		cdPath = absPath
		return []byte{}
	}

	fmt.Printf("cd: %v: is not a directory\n", absPath)
	return []byte{}
}

func getPathString(tokens []token.Token) string {
	var pathString string = ""
	for _, tok := range tokens {
		pathString += tok.Literal
	}

	return pathString
}
