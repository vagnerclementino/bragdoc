package handler

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CmdHandler interface {
	Register(cmd string) error
	Execute() error
}

type cmdHandler struct {
	root *cobra.Command
	cmds []string
}

func (h *cmdHandler) Register(cmd string) error {
	var err error

	switch cmd {
	case "version":
		h.root.AddCommand(versionCmd)
	default:
		err = fmt.Errorf("the command '%s' cannot be registered", cmd)
	}
	return err
}

func (h *cmdHandler) Execute() error {
	return h.root.Execute()
}

func NewCmdHandler() CmdHandler {
	return &cmdHandler{
		root: rootCmd,
		cmds: []string{},
	}
}
