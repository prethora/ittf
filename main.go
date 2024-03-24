package main

import (
	"fmt"
	"os"

	"github.com/prethora/ittf/internal/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
