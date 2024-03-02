package cmd

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
	Short: `Display the current version and build of the bradoc`,
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

		if _, err := fmt.Fprintln(writer, "Version:\t", Version); err != nil {
			return
		}
		if _, err := fmt.Fprintln(writer, "Build:\t", Build); err != nil {
			return
		}
	},
}
