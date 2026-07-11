package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/myerscode/aws-meta/examples/shared"
	"github.com/myerscode/aws-meta/pkg/partitions"
)

func main() {
	log.Println("=== AWS Partitions Detailed Information ===")

	partitionList, err := partitions.List()
	if err != nil {
		log.Fatalf("Failed to get partition list: %v", err)
	}

	log.Printf("Found %d partitions:", len(partitionList))
	fmt.Println()

	for _, partition := range partitionList {
		fmt.Printf("Partition: %s\n", partition.ID)
		fmt.Printf("  DNS Suffix: %s\n", partition.DNSSuffix)
		fmt.Printf("  Dual Stack DNS: %s\n", partition.DualStackDNSSuffix)
		fmt.Printf("  Region Regex: %s\n", partition.RegionRegex)
		fmt.Printf("  Global Region: %s\n", partition.ImplicitGlobalRegion)
		fmt.Printf("  Regions: %d\n", len(partition.Regions))

		regionCount := shared.Min(5, len(partition.Regions))
		if partition.ID == "aws" {
			regionCount = shared.Min(8, len(partition.Regions))
		}

		if regionCount > 0 {
			for i := 0; i < regionCount; i++ {
				region := partition.Regions[i]
				fmt.Printf("    - %s (%s)\n", region.RegionId, region.RegionName)
			}
			if len(partition.Regions) > regionCount {
				fmt.Printf("    ... and %d more\n", len(partition.Regions)-regionCount)
			}
		}
		fmt.Println()
	}

	// Summary by partition type
	log.Println("=== Partition Summary ===")
	commercial, sovereign, isolated := 0, 0, 0
	for _, partition := range partitionList {
		switch {
		case partition.ID == "aws":
			commercial++
		case strings.Contains(partition.ID, "iso"):
			isolated++
		default:
			sovereign++
		}
	}
	fmt.Printf("Commercial: %d, Sovereign: %d, Isolated: %d\n", commercial, sovereign, isolated)
}
