package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		latestInput := strings.Split(input[:len(input)-1], " ")
		command := latestInput[0]
		args := latestInput[1:]

		switch command {
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "exit":
			os.Exit(0)
		default:
			fmt.Println(command + ": command not found")
		}
	}
}
