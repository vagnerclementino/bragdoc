package main

import (
	"github.com/vagnerclementino/bragdoc/internal/handler/cmd"
	"os"
)

func main() {
	cli, err := cmd.NewCli([]string{
		"version",
	})
	if err != nil {
		os.Exit(1)
	}

	if err := cli.Start(); err != nil {
		os.Exit(1)
	}
}
