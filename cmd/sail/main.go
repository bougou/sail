package main

import (
	"fmt"
	"os"

	"github.com/bougou/sail/pkg/commands"
)

func main() {
	rootCmd := commands.NewSailCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
