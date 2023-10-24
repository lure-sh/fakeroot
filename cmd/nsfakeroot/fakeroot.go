package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"lure.sh/fakeroot"
	"lure.sh/fakeroot/loginshell"
)

func main() {
	var showHelp bool
	flag.BoolVar(&showHelp, "help", false, "Show help screen")
	flag.BoolVar(&showHelp, "h", false, "Show help screen")
	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	var (
		cmd  string
		args []string
		err  error
	)
	if flag.NArg() > 0 {
		cmd = flag.Arg(0)
		args = flag.Args()[1:]
	} else {
		cmd, err = loginshell.Get(-1)
		if err != nil {
			log.Fatalln(err)
		}
	}

	c, err := fakeroot.Command(cmd, args...)
	if err != nil {
		log.Fatalln(err)
	}

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	err = c.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func printHelp() {
	fmt.Print("Go implementation of fakeroot using Linux user namespaces.\n\n")
	fmt.Print("Usage: fakeroot [cmd] [args...]\n\n")
	fmt.Print("Arguments:\n")
	fmt.Print(" [cmd]     Command to execute in fakeroot environment. If not specified, the user's login shell will be executed.\n")
	fmt.Print(" [args...] Arguments to pass to the executed command.\n")
}
