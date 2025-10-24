package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/myerscode/aws-meta/internal/aws"
	"github.com/myerscode/aws-meta/pkg/partitions"
	"github.com/spf13/cobra"
)

// listPartitionsCmd represents the list partitions command
var listPartitionsCmd = &cobra.Command{
	Use:   "partitions",
	Short: "List AWS Partition information",
	Long: `List AWS Partition metadata including Partition IDs, DNS suffixes,
region regex patterns, and the list of regions within each Partition.

Use the --json flag to output the data in JSON format for programmatic use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOutput, _ := cmd.Flags().GetBool("json")

		if jsonOutput {
			// Get Partition data using existing pkg function
			partitionData, err := partitions.List()
			if err != nil {
				return fmt.Errorf("failed to retrieve Partition data: %w", err)
			}

			// Handle empty data case - return empty array if no data available
			if partitionData == nil {
				partitionData = make(aws.PartitionSchemas, 0)
			}

			// Marshal to JSON with proper indentation
			jsonBytes, err := json.MarshalIndent(partitionData, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal Partition data to JSON: %w", err)
			}

			fmt.Println(string(jsonBytes))
		} else {
			// Todo placeholder for non-JSON mode
			fmt.Println("TODO: Human-readable Partition listing not yet implemented. Use --json flag for JSON output.")
		}

		return nil
	},
}

func init() {
	// Add partitions subcommand to list command
	listCmd.AddCommand(listPartitionsCmd)

	// Add --json flag with help text
	listPartitionsCmd.Flags().BoolP("json", "j", false, "Output partition data in JSON format")
}
