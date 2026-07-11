package commands

import (
	"encoding/json"
	"fmt"

	"github.com/myerscode/aws-meta/internal/aws"
	"github.com/myerscode/aws-meta/pkg/services"
	"github.com/spf13/cobra"
)

var listServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List AWS Service information",
	Long: `List AWS Service metadata including Service IDs, full names,
operations, and regional availability information.

Use the --json flag to output the data in JSON format for programmatic use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOutput, _ := cmd.Flags().GetBool("json")

		if jsonOutput {
			serviceData, err := services.List()
			if err != nil {
				return fmt.Errorf("failed to retrieve Service data: %w", err)
			}

			if serviceData == nil {
				serviceData = make(aws.ServiceSchemas)
			}

			jsonBytes, err := json.MarshalIndent(serviceData, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal Service data to JSON: %w", err)
			}

			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println("TODO: Human-readable Service listing not yet implemented. Use --json flag for JSON output.")
		}

		return nil
	},
}

func init() {
	listCmd.AddCommand(listServicesCmd)
	listServicesCmd.Flags().BoolP("json", "j", false, "Output service data in JSON format")
}
