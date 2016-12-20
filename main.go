package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// builtinCmds the commands this shell supports.
var builtinCmds = []string{
	"cd",
	"help",
	"exit",
}

// cmdLoop main loop of the shell, interprets and executes commands.
func cmdLoop() {
	var status int
	for status == 0 {
		fmt.Printf("> ")
		line := readLine()
		args := splitLine(line)
		status = execute(args)
	}
}

// readLine reads input up to newline from stdin.
func readLine() string {
	r := bufio.NewReader(os.Stdin)
	l, err := r.ReadString('\n')
	if err != nil {
		os.Exit(1)
	}
	return strings.TrimRight(l, "\n")
}

// splitLine splits line into individual strings.
func splitLine(line string) []string {
	if len(line) == 0 {
		return []string{}
	}
	// if user enters command then space e.g. 'ls[SPACE]' space is interpreted
	// as an argument to ls, remove it.
	var args []string
	for _, str := range strings.Split(line, " ") {
		if len(str) > 0 {
			args = append(args, strings.TrimSpace(str))
		}
	}
	return args
}

// execute determines which command to execute.
func execute(args []string) int {
	if len(args) == 0 {
		return 0
	}
	var f func([]string) int
	switch args[0] {
	case "cd":
		f = baresh_cd
	case "help":
		f = baresh_help
	case "exit":
		f = baresh_exit
	default:
		f = launch
	}
	return f(args)
}

// launch runs external commands ls, tar, cat, etc...
func launch(args []string) int {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = strings.NewReader("")
	var serr bytes.Buffer
	cmd.Stderr = &serr
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		if serr.Len() > 0 {
			// external command returned status > 0
			fmt.Print(serr.String())
		} else {
			// command not found in $PATH
			fmt.Println("-baresh:", args[0], "command not found")
		}
	}
	if out.Len() > 0 {
		fmt.Print(out.String())
	}
	return 0
}

// internal commands

// baresh_cd change process working directory.
func baresh_cd(args []string) int {
	if len(args) == 1 {
		// no directory name passed, should cd to $HOME. we'll just return :).
		return 0
	}
	err := os.Chdir(args[1])
	if err != nil {
		fmt.Println("-baresh:", args[1], "no such file or directory")
	}
	return 0
}

// baresh_help print help information.
func baresh_help(args []string) int {
	const help = `baresh 

These shell commands are defined internally. Type 'help' to see this list.

help
cd <directory name>
exit
`
	fmt.Print(help)
	return 0
}

// baresh_exit terminate the process.
func baresh_exit(args []string) int {
	// break out of the shells loop
	return 1
}

func main() {
	// 1. Load config files, if any.

	// 2. Run command loop.
	cmdLoop()

	// 3. Perform any shutdown/cleanup.
}
