package main

import (
	"fmt"

	"github.com/myerscode/aws-meta/internal/util"
	"github.com/myerscode/aws-meta/pkg/partitions"
)

func main() {
	// Test Commercial Partitions
	util.LogInfo("=== Commercial Partitions ===")
	commercial, err := partitions.CommercialPartitions()
	if err != nil {
		util.LogError(fmt.Sprintf("Failed to get commercial partitions: %v", err))
	} else {
		util.LogInfo(fmt.Sprintf("Found %d commercial partitions:", len(commercial)))
		for _, partition := range commercial {
			util.LogInfo(fmt.Sprintf("  - %s", partition))
		}
	}
	fmt.Println()

	// Test Sovereign Partitions
	util.LogInfo("=== Sovereign Partitions ===")
	sovereign, err := partitions.SovereignPartitions()
	if err != nil {
		util.LogError(fmt.Sprintf("Failed to get sovereign partitions: %v", err))
	} else {
		util.LogInfo(fmt.Sprintf("Found %d sovereign partitions:", len(sovereign)))
		for _, partition := range sovereign {
			util.LogInfo(fmt.Sprintf("  - %s", partition))
		}
	}
	fmt.Println()

	// Test Isolated Partitions
	util.LogInfo("=== Isolated Partitions ===")
	isolated, err := partitions.IsolatedPartitions()
	if err != nil {
		util.LogError(fmt.Sprintf("Failed to get isolated partitions: %v", err))
	} else {
		util.LogInfo(fmt.Sprintf("Found %d isolated partitions:", len(isolated)))
		for _, partition := range isolated {
			util.LogInfo(fmt.Sprintf("  - %s", partition))
		}
	}
	fmt.Println()

	// Summary
	util.LogInfo("=== Summary ===")
	allPartitions := partitions.AllPartitionNames()
	totalCategorized := len(commercial) + len(sovereign) + len(isolated)

	util.LogInfo(fmt.Sprintf("Total partitions: %d", len(allPartitions)))
	util.LogInfo(fmt.Sprintf("Commercial: %d", len(commercial)))
	util.LogInfo(fmt.Sprintf("Sovereign: %d", len(sovereign)))
	util.LogInfo(fmt.Sprintf("Isolated: %d", len(isolated)))
	util.LogInfo(fmt.Sprintf("Categorized: %d", totalCategorized))

	if totalCategorized == len(allPartitions) {
		util.LogInfo("✓ All partitions are properly categorized")
	} else {
		util.LogError(fmt.Sprintf("✗ Mismatch: %d total vs %d categorized", len(allPartitions), totalCategorized))
	}

	// Show partition characteristics
	util.LogInfo("=== Partition Characteristics ===")
	util.LogInfo("Commercial partitions:")
	util.LogInfo("  - Standard AWS partition available to general public")
	util.LogInfo("  - DNS suffix: amazonaws.com")

	util.LogInfo("Sovereign partitions:")
	util.LogInfo("  - Operated by local entities for regulatory compliance")
	util.LogInfo("  - aws-cn: China partition (amazonaws.com.cn)")
	util.LogInfo("  - aws-us-gov: US Government partition (amazonaws.com)")
	util.LogInfo("  - aws-eusc: EU Sovereign Cloud (amazonaws.eu)")

	util.LogInfo("Isolated partitions:")
	util.LogInfo("  - Air-gapped environments for highly sensitive workloads")
	util.LogInfo("  - Various government and intelligence community partitions")
	util.LogInfo("  - Different DNS suffixes for security isolation")
}
