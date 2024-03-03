package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/handler"
)

type cmdHandler struct {
	root *cobra.Command
}

func (h *cmdHandler) Start() error {
	return h.root.Execute()
}

func NewCli(cmds []string) (handler.Handler, error) {
	cli := &cmdHandler{
		root: rootCmd,
	}
	for _, cmd := range cmds {
		if err := cli.register(cmd); err != nil {
			return nil, err
		}
	}
	return cli, nil
}

func (h *cmdHandler) register(cmd string) error {
	var err error
	switch cmd {
	case "version":
		h.root.AddCommand(versionCmd)
	case "init":
		h.root.AddCommand(initCmd)
	case "brag":
		h.root.AddCommand(bragCmd)
	case "doc":
		h.root.AddCommand(docCmd)
	default:
		err = fmt.Errorf("the command '%s' cannot be registered", cmd)
	}
	return err
}
