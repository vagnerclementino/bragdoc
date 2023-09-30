package main

import (
	"github.com/vagnerclementino/bragdoc/internal/handler"
	"os"
)

func main() {
	cmdHandler := handler.NewCmdHandler()
	cmdHandler.CmdRegister("version")
	if err := cmdHandler.Execute(); err != nil {
		os.Exit(1)
	}
}
