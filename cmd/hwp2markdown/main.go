package main

import (
	"os"

	"github.com/roboco-io/hwp2markdown/internal/cli"
)

var version = "dev"

func main() {
	cli.SetVersion(version)
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
