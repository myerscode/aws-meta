package commands

import (
	"encoding/json"
	"fmt"

	"github.com/myerscode/aws-meta/pkg/regions"
	"github.com/spf13/cobra"
)

var listRegionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "List AWS Region information",
	Long: `List AWS Region metadata including Region IDs, names, partition information,
and the list of services available in each region.

Use the --json flag to output the data in JSON format for programmatic use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOutput, _ := cmd.Flags().GetBool("json")

		if jsonOutput {
			regionData, err := regions.ListAllRegions()
			if err != nil {
				return fmt.Errorf("failed to retrieve Region data: %w", err)
			}

			if regionData == nil {
				regionData = make([]regions.RegionInfo, 0)
			}

			jsonBytes, err := json.MarshalIndent(regionData, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal Region data to JSON: %w", err)
			}

			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println("TODO: Human-readable Region listing not yet implemented. Use --json flag for JSON output.")
		}

		return nil
	},
}

func init() {
	listCmd.AddCommand(listRegionsCmd)
	listRegionsCmd.Flags().BoolP("json", "j", false, "Output Region data in JSON format")
}
