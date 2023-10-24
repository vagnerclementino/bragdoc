package main

import (
	"os"

	"github.com/vagnerclementino/bragdoc/internal/handler"
)

func main() {
	cmdHandler := handler.NewCmdHandler()
	cmdHandler.Register("version")
	if err := cmdHandler.Execute(); err != nil {
		os.Exit(1)
	}
}
