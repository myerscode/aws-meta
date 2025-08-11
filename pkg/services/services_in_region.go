package services

import (
	"fmt"
)

// ServicesInRegion returns a list of service names available in the specified region.
// The returned services are sorted alphabetically.
//
// Parameters:
//   - regionName: The AWS region name (e.g., "us-east-1", "eu-west-1")
//
// Returns:
//   - []string: A sorted list of service names available in the region
//   - error: An error if the region name is invalid or if there's an issue reading the manifest
//
// Example:
//
//	services, err := ServicesInRegion("us-east-1")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d services in us-east-1\n", len(services))
func ServicesInRegion(regionName string) ([]string, error) {
	manifest, err := regionManifest()
	if err != nil {
		return nil, err
	}

	// Search through all partitions and regions to find the specified region
	for _, partition := range manifest {
		for _, region := range partition.Regions {
			if region.RegionName == regionName {
				return region.Services, nil
			}
		}
	}

	// If we get here, the region was not found
	return nil, fmt.Errorf("invalid region name: %s", regionName)
}
