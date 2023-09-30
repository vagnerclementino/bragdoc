/*
Copyright © 2023 Vagner Clementino vagner.clemetino@gmail.com
*/
package handler

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	Version = "unknown"
	Build   = "unknown"
)

const defaultWidthTable = 8

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: `ℹ️  This command is used to display the current version of the application`,
	Long:  `When executed, it will provide you with information about the version number of the application you're running. 🚀`,
	Run: func(cmd *cobra.Command, args []string) {
		writer := new(tabwriter.Writer)

		writer.Init(os.Stdout, defaultWidthTable, defaultWidthTable, 0, '\t', 0)

		defer func(writer *tabwriter.Writer) {
			err := writer.Flush()
			if err != nil {
				return
			}
		}(writer)

		fmt.Fprintln(writer, "Version:\t", Version)
		fmt.Fprintln(writer, "Build:\t", Build)

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
