/*
Copyright © 2023 Vagner Clementino vagner.clementino@gmail.com
*/
package handler

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bragdoc",
	Short: "A brief description of your application",
	Long:  `Bragdoc is a powerful command-line interface (CLI) tool designed to help individuals build their own "Brag Documents."`,
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bragdoc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
