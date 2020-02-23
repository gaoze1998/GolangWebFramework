package main

import (
	"github.com/gaoze1998/GolangWebFramework/Cli"
	"os"
)

func main() {
	args := os.Args
	command := args[1]
	switch command {
	case "create":
		Cli.Create(args)

	}
}
