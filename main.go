package main

import (
	"os"

	"github.com/rduarte/skl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
