package commands

import (
	"github.com/myerscode/aws-meta/cli/internal/util"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List AWS metadata information",
	Long: `List various AWS metadata including partitions and services.

Use the available subcommands to list specific types of metadata:
  - partitions: List AWS partition information
  - services: List AWS service information
  - regions: List AWS region information

Each subcommand supports optional flags for different output formats.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			util.LogError("Error displaying help: " + err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
