package main

import (
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

func main() {
	var userCommand string
	fmt.Fprint(os.Stdout, "$ ")
	fmt.Scanln(&userCommand)
	msg := fmt.Sprintf("%s: command not found", userCommand)
	fmt.Fprint(os.Stdout, msg)
}
