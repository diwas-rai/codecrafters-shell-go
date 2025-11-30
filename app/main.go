package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

var COMMAND_WORDS = []string{"echo", "exit", "type", "pwd", "cd"}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		argv := parseInput(input)
		if len(argv) == 0 {
			continue
		}
		command := argv[0]

		switch command {
		case "echo":
			echoCommand(argv)
		case "exit":
			exitCommand(argv)
		case "type":
			typeCommand(argv)
		case "pwd":
			pwdCommand()
		case "cd":
			cdCommand(argv)
		default:
			execute(argv)
		}
	}
}

func echoCommand(argv []string) {
	fmt.Println(strings.Join(argv[1:], " "))
}

func exitCommand(argv []string) {
	code := 0

	if len(argv) > 1 {
		argCode, err := strconv.Atoi(argv[1])
		if err == nil {
			code = argCode
		}
	}

	os.Exit(code)
}

func typeCommand(argv []string) {
	if len(argv) == 1 {
		return
	}

	val := argv[1]

	if slices.Contains(COMMAND_WORDS, val) {
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", val)
		return
	}

	if file, exists := findBinInPath(val); exists {
		fmt.Fprintf(os.Stdout, "%s is %s\n", val, file)
		return
	}

	fmt.Fprintf(os.Stdout, "%s: not found\n", val)
}

func execute(argv []string) {
	if len(argv) == 0 {
		return
	}

	fileName := argv[0]

	if _, exists := findBinInPath(fileName); exists {
		cmd := exec.Command(fileName, argv[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err == nil {
			return
		}
	}

	fmt.Fprintf(os.Stdout, "%s: command not found\n", fileName)
}

func pwdCommand() {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stdout, "%s\n", wd)
}

func cdCommand(argv []string) {
	if len(argv) == 1 {
		return
	}

	path := argv[1]
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", path)
		return
	}

	path = strings.Replace(path, "~", home, 1)

	err = os.Chdir(path)
	if err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", path)
	}
}

func findBinInPath(bin string) (string, bool) {
	paths := os.Getenv("PATH")

	for _, path := range strings.Split(paths, string(os.PathListSeparator)) {
		file := filepath.Join(path, bin)
		fileInfo, err := os.Stat(file)

		if err == nil && !fileInfo.IsDir() && fileInfo.Mode().Perm()&0111 != 0 {
			return file, true
		}
	}

	return "", false
}

func parseInput(input string) []string {
	var args []string
	var currentArg strings.Builder

	inSingleQuote := false
	inDoubleQuote := false

	for _, r := range input {
		switch r {
		case '\n':
			continue
		case '\'':
			if inDoubleQuote {
				currentArg.WriteRune(r)
			} else {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if inSingleQuote {
				currentArg.WriteRune(r)
			} else {
				inDoubleQuote = !inDoubleQuote
			}
		case ' ':
			if !inSingleQuote && !inDoubleQuote {
				if currentArg.Len() > 0 {
					args = append(args, currentArg.String())
					currentArg.Reset()
				}
			} else {
				currentArg.WriteRune(r)
			}
		default:
			currentArg.WriteRune(r)
		}
	}

	// If there is a leftover argument at the end, add it
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}
