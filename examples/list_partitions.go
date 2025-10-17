package main

import (
	"fmt"

	"github.com/myerscode/aws-meta/internal/util"
	"github.com/myerscode/aws-meta/pkg/partitions"
)

func main() {
	util.LogInfo("=== AWS Partitions Detailed Information ===")

	partitionList, err := partitions.ListPartitions()
	if err != nil {
		util.LogError(fmt.Sprintf("Failed to get partition list: %v", err))
		return
	}

	util.LogInfo(fmt.Sprintf("Found %d partitions:", len(partitionList)))
	fmt.Println()

	for _, partition := range partitionList {
		util.LogInfo(fmt.Sprintf("Partition: %s", partition.ID))
		util.LogInfo(fmt.Sprintf("  DNS Suffix: %s", partition.DNSSuffix))
		util.LogInfo(fmt.Sprintf("  Dual Stack DNS: %s", partition.DualStackDNSSuffix))
		util.LogInfo(fmt.Sprintf("  Region Regex: %s", partition.RegionRegex))
		util.LogInfo(fmt.Sprintf("  Global Region: %s", partition.ImplicitGlobalRegion))
		util.LogInfo(fmt.Sprintf("  Regions: %d", len(partition.Regions)))

		// Show regions (more for AWS partition, fewer for others)
		var regionCount int
		if partition.ID == "aws" {
			regionCount = min(8, len(partition.Regions)) // Show more for AWS
		} else {
			regionCount = min(5, len(partition.Regions)) // Show fewer for others
		}

		if regionCount > 0 {
			if len(partition.Regions) <= regionCount {
				util.LogInfo("  All regions:")
			} else {
				util.LogInfo("  Sample regions:")
			}

			for i := 0; i < regionCount; i++ {
				region := partition.Regions[i]
				util.LogInfo(fmt.Sprintf("    - %s (%s)", region.RegionId, region.RegionName))
			}
			if len(partition.Regions) > regionCount {
				util.LogInfo(fmt.Sprintf("    ... and %d more", len(partition.Regions)-regionCount))
			}
		}
		fmt.Println()
	}

	// Summary by partition type
	util.LogInfo("=== Partition Summary ===")
	commercial := 0
	sovereign := 0
	isolated := 0

	for _, partition := range partitionList {
		switch {
		case partition.ID == "aws":
			commercial++
		case containsString(partition.ID, "iso"):
			isolated++
		default:
			sovereign++
		}
	}

	util.LogInfo(fmt.Sprintf("Commercial partitions: %d", commercial))
	util.LogInfo(fmt.Sprintf("Sovereign partitions: %d", sovereign))
	util.LogInfo(fmt.Sprintf("Isolated partitions: %d", isolated))

	// Show DNS suffix patterns
	util.LogInfo("=== DNS Suffix Patterns ===")
	dnsSuffixes := make(map[string][]string)
	for _, partition := range partitionList {
		suffix := partition.DNSSuffix
		dnsSuffixes[suffix] = append(dnsSuffixes[suffix], partition.ID)
	}

	for suffix, partitionIds := range dnsSuffixes {
		util.LogInfo(fmt.Sprintf("%s: %v", suffix, partitionIds))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					hasSubstring(s, substr))))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
