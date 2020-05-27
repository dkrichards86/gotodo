package main

import (
	"os"

	"github.com/dkrichards86/gotodo/internal/commands"
)

func main() {
	commands.RunCli(os.Args)
}
