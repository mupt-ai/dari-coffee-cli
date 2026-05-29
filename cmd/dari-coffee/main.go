package main

import (
	"os"

	"github.com/mupt-ai/dari-coffee-cli/internal/cli"
)

var version = "dev"

func main() {
	if err := cli.Execute(version); err != nil {
		os.Exit(1)
	}
}
