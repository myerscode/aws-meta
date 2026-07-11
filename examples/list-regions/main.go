package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/myerscode/aws-meta/pkg/regions"
)

func main() {
	log.Println("=== AWS Regions Detailed Information ===")

	regionList, err := regions.ListAllRegions()
	if err != nil {
		log.Fatalf("Failed to get region list: %v", err)
	}

	log.Printf("Found %d regions across all partitions:", len(regionList))
	fmt.Println()

	// Group regions by partition
	partitionRegions := make(map[string][]regions.RegionInfo)
	for _, region := range regionList {
		partitionRegions[region.PartitionID] = append(partitionRegions[region.PartitionID], region)
	}

	var partitionNames []string
	for partition := range partitionRegions {
		partitionNames = append(partitionNames, partition)
	}
	sort.Strings(partitionNames)

	for _, partitionName := range partitionNames {
		regs := partitionRegions[partitionName]
		fmt.Printf("Partition: %s (%d regions)\n", partitionName, len(regs))

		displayCount := min(8, len(regs))
		if partitionName != "aws" {
			displayCount = min(5, len(regs))
		}

		for i := 0; i < displayCount; i++ {
			region := regs[i]
			fmt.Printf("  - %s (%s) - %d services\n",
				region.RegionId, region.RegionName, len(region.Services))
		}

		if len(regs) > displayCount {
			fmt.Printf("  ... and %d more regions\n", len(regs)-displayCount)
		}
		fmt.Println()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
