package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1
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

		var out io.Writer = os.Stdout

		for i, arg := range argv {
			if arg == ">" || arg == "1>" {
				if i+1 < len(argv) {
					filePath := argv[i+1]

					// Open file: Create if missing, Write Only, Truncate content
					f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
						argv = nil
						break
					}
					out = f

					argv = append(argv[:i], argv[i+2:]...)
				}
				break
			}
		}

		// If argv became empty or invalid due to redirection error
		if len(argv) == 0 {
			continue
		}

		command := argv[0]

		switch command {
		case "echo":
			echoCommand(argv, out)
		case "exit":
			exitCommand(argv)
		case "type":
			typeCommand(argv, out)
		case "pwd":
			pwdCommand(out)
		case "cd":
			cdCommand(argv, out)
		default:
			execute(argv, out)
		}
	}
}

func echoCommand(argv []string, out io.Writer) {
	fmt.Fprintln(out, strings.Join(argv[1:], " "))
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

func typeCommand(argv []string, out io.Writer) {
	if len(argv) == 1 {
		return
	}

	val := argv[1]

	if slices.Contains(COMMAND_WORDS, val) {
		fmt.Fprintf(out, "%s is a shell builtin\n", val)
		return
	}

	if file, exists := findBinInPath(val); exists {
		fmt.Fprintf(out, "%s is %s\n", val, file)
		return
	}

	fmt.Fprintf(out, "%s: not found\n", val)
}

func execute(argv []string, out io.Writer) {
	if len(argv) == 0 {
		return
	}

	fileName := argv[0]

	if _, exists := findBinInPath(fileName); exists {
		cmd := exec.Command(fileName, argv[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = out

		cmd.Run()
		return
	}

	fmt.Fprintf(out, "%s: command not found\n", fileName)
}

func pwdCommand(out io.Writer) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Fprintf(out, "%s\n", wd)
}

func cdCommand(argv []string, out io.Writer) {
	if len(argv) == 1 {
		return
	}

	path := argv[1]
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(out, "cd: %s: No such file or directory\n", path)
		return
	}

	path = strings.Replace(path, "~", home, 1)

	err = os.Chdir(path)
	if err != nil {
		fmt.Fprintf(out, "cd: %s: No such file or directory\n", path)
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
	runes := []rune(input)
	var args []string
	var sb strings.Builder

	inSingle, inDouble := false, false

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}
		if r == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}

		if r == '\\' {
			if inSingle {
				sb.WriteRune(r)
				continue
			}
			if i+1 < len(runes) {
				next := runes[i+1]
				if inDouble {
					if strings.ContainsRune("$`\"\\\n", next) {
						sb.WriteRune(next)
						i++
					} else {
						sb.WriteRune(r)
					}
				} else {
					sb.WriteRune(next)
					i++
				}
				continue
			}
		}

		if r == ' ' && !inSingle && !inDouble {
			if sb.Len() > 0 {
				args = append(args, sb.String())
				sb.Reset()
			}
			continue
		}

		if r == '\n' {
			continue
		}

		sb.WriteRune(r)
	}

	if sb.Len() > 0 {
		args = append(args, sb.String())
	}

	return args
}
