package regions

import (
	"github.com/myerscode/aws-meta/internal/aws"
	"github.com/myerscode/aws-meta/pkg/data"
)

// RegionInfo represents detailed information about an AWS region
type RegionInfo struct {
	RegionId    string   `json:"regionId"`
	RegionName  string   `json:"regionName"`
	PartitionID string   `json:"partitionId"`
	Services    []string `json:"services"`
}

// partitionManifest loads the partition manifest data
func partitionManifest() (aws.PartitionSchemas, error) {
	manifest, err := data.PartitionManifest()
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

// regionManifest loads the region manifest data
func regionManifest() (aws.RegionSchemas, error) {
	manifest, err := data.RegionsManifest()
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

// ListAllRegions returns detailed information about all AWS regions across all partitions.
// This includes region metadata, partition information, and the list of services available in each region.
//
// Returns:
//   - []RegionInfo: A slice of region information containing detailed region data
//   - error: An error if there's an issue reading the manifest
//
// Example:
//
//	regions, err := regions.ListAllRegions()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, region := range regions {
//	    fmt.Printf("Region: %s (%s)\n", region.RegionId, region.RegionName)
//	    fmt.Printf("  Partition: %s\n", region.PartitionID)
//	    fmt.Printf("  Services: %d\n", len(region.Services))
//	}
func ListAllRegions() ([]RegionInfo, error) {
	// Get partition data for region names and partition mapping
	partitions, err := partitionManifest()
	if err != nil {
		return nil, err
	}

	// Get region data for service information
	regionSchemas, err := regionManifest()
	if err != nil {
		return nil, err
	}

	// Create a map of region ID to services for quick lookup
	regionServices := make(map[string][]string)
	for _, schema := range regionSchemas {
		for _, region := range schema.Regions {
			regionServices[region.RegionName] = region.Services
		}
	}

	// Build the comprehensive region list
	var allRegions []RegionInfo
	for _, partition := range partitions {
		for _, region := range partition.Regions {
			services := regionServices[region.RegionId]
			if services == nil {
				services = []string{} // Empty slice if no services found
			}

			regionInfo := RegionInfo{
				RegionId:    region.RegionId,
				RegionName:  region.RegionName,
				PartitionID: partition.ID,
				Services:    services,
			}
			allRegions = append(allRegions, regionInfo)
		}
	}

	return allRegions, nil
}
