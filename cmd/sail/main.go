package main

import (
	"fmt"
	"os"

	"github.com/bougou/sail/pkg/cmd"
)

func main() {
	rootCmd := cmd.NewSailCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
