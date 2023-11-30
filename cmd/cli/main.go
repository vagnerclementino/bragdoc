package main

import (
	"os"

	"github.com/vagnerclementino/bragdoc/internal/handler"
)

func main() {
	cmdHandler := handler.NewCmdHandler()
	if err := cmdHandler.Register("version"); err != nil {
		os.Exit(1)
	}
	if err := cmdHandler.Execute(); err != nil {
		os.Exit(1)
	}
}
