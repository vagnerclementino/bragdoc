package handler

import "github.com/spf13/cobra"

type CmdHandler interface {
	CmdRegister(cmd string)
	Execute() error
}

type cmdHandler struct {
	root *cobra.Command
	cmds []string
}

func (h *cmdHandler) CmdRegister(cmd string) {
	h.cmds = append(h.cmds, cmd)
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
