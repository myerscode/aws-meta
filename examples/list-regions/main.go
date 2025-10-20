package main

import (
	"fmt"
	"sort"

	"github.com/myerscode/aws-meta/internal/util"
	"github.com/myerscode/aws-meta/pkg/regions"
)

func main() {
	util.LogInfo("=== AWS Regions Detailed Information ===")

	regionList, err := regions.ListAllRegions()
	if err != nil {
		util.LogError(fmt.Sprintf("Failed to get region list: %v", err))
		return
	}

	util.LogInfo(fmt.Sprintf("Found %d regions across all partitions:", len(regionList)))
	fmt.Println()

	// Group regions by partition for better organization
	partitionRegions := make(map[string][]regions.RegionInfo)
	for _, region := range regionList {
		partitionRegions[region.PartitionID] = append(partitionRegions[region.PartitionID], region)
	}

	// Sort partition names for consistent output
	var partitionNames []string
	for partition := range partitionRegions {
		partitionNames = append(partitionNames, partition)
	}
	sort.Strings(partitionNames)

	// Display regions grouped by partition
	for _, partitionName := range partitionNames {
		regions := partitionRegions[partitionName]
		util.LogInfo(fmt.Sprintf("Partition: %s (%d regions)", partitionName, len(regions)))

		// Show first several regions for each partition
		displayCount := min(8, len(regions))
		if partitionName != "aws" {
			displayCount = min(5, len(regions)) // Show fewer for non-AWS partitions
		}

		for i := 0; i < displayCount; i++ {
			region := regions[i]
			serviceCount := len(region.Services)
			util.LogInfo(fmt.Sprintf("  - %s (%s) - %d services",
				region.RegionId, region.RegionName, serviceCount))
		}

		if len(regions) > displayCount {
			util.LogInfo(fmt.Sprintf("  ... and %d more regions", len(regions)-displayCount))
		}
		fmt.Println()
	}

	// Summary statistics
	util.LogInfo("=== Regional Statistics ===")

	// Count regions by partition
	for _, partitionName := range partitionNames {
		regions := partitionRegions[partitionName]
		util.LogInfo(fmt.Sprintf("%s: %d regions", partitionName, len(regions)))
	}
	fmt.Println()

	// Find regions with most and least services
	var maxServices, minServices = 0, 999999
	var maxRegion, minRegion regions.RegionInfo

	for _, region := range regionList {
		serviceCount := len(region.Services)
		if serviceCount > maxServices {
			maxServices = serviceCount
			maxRegion = region
		}
		if serviceCount < minServices && serviceCount > 0 { // Ignore regions with 0 services
			minServices = serviceCount
			minRegion = region
		}
	}

	if maxServices > 0 {
		util.LogInfo(fmt.Sprintf("Region with most services: %s (%s) - %d services",
			maxRegion.RegionId, maxRegion.RegionName, maxServices))
	}
	if minServices < 999999 {
		util.LogInfo(fmt.Sprintf("Region with least services: %s (%s) - %d services",
			minRegion.RegionId, minRegion.RegionName, minServices))
	}

	// Show geographic distribution
	util.LogInfo("=== Geographic Distribution (AWS Partition) ===")
	awsRegions := partitionRegions["aws"]
	if len(awsRegions) > 0 {
		geographicAreas := make(map[string][]string)

		for _, region := range awsRegions {
			// Simple geographic grouping based on region prefix
			switch region.RegionId[:2] {
			case "us":
				geographicAreas["North America"] = append(geographicAreas["North America"], region.RegionId)
			case "eu":
				geographicAreas["Europe"] = append(geographicAreas["Europe"], region.RegionId)
			case "ap":
				geographicAreas["Asia Pacific"] = append(geographicAreas["Asia Pacific"], region.RegionId)
			case "sa":
				geographicAreas["South America"] = append(geographicAreas["South America"], region.RegionId)
			case "ca":
				geographicAreas["Canada"] = append(geographicAreas["Canada"], region.RegionId)
			case "af":
				geographicAreas["Africa"] = append(geographicAreas["Africa"], region.RegionId)
			case "me":
				geographicAreas["Middle East"] = append(geographicAreas["Middle East"], region.RegionId)
			case "il":
				geographicAreas["Israel"] = append(geographicAreas["Israel"], region.RegionId)
			case "mx":
				geographicAreas["Mexico"] = append(geographicAreas["Mexico"], region.RegionId)
			default:
				geographicAreas["Other"] = append(geographicAreas["Other"], region.RegionId)
			}
		}

		for area, regionIds := range geographicAreas {
			util.LogInfo(fmt.Sprintf("%s: %d regions (%v)", area, len(regionIds), regionIds[:min(3, len(regionIds))]))
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
