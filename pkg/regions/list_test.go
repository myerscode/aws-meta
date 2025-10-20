package regions

import (
	"testing"
)

func TestListAllRegions(t *testing.T) {
	regions, err := ListAllRegions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(regions) == 0 {
		t.Fatal("Expected at least one region")
	}

	// Check that all regions have required fields
	for i, region := range regions {
		if region.RegionId == "" {
			t.Errorf("Region at index %d has empty RegionId", i)
		}
		if region.RegionName == "" {
			t.Errorf("Region %s has empty RegionName", region.RegionId)
		}
		if region.PartitionID == "" {
			t.Errorf("Region %s has empty PartitionID", region.RegionId)
		}
		// Services can be empty for some regions, so we don't check for that
	}
}

func TestListAllRegionsContent(t *testing.T) {
	regions, err := ListAllRegions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check for known regions
	expectedRegions := []string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-1", "cn-north-1", "us-gov-west-1"}
	foundRegions := make(map[string]bool)

	for _, region := range regions {
		foundRegions[region.RegionId] = true
	}

	for _, expected := range expectedRegions {
		if !foundRegions[expected] {
			t.Errorf("Expected region %s not found", expected)
		}
	}

	// Check that regions are from different partitions
	partitions := make(map[string]bool)
	for _, region := range regions {
		partitions[region.PartitionID] = true
	}

	expectedPartitions := []string{"aws", "aws-cn", "aws-us-gov"}
	for _, expected := range expectedPartitions {
		if !partitions[expected] {
			t.Errorf("Expected partition %s not found in regions", expected)
		}
	}
}

func TestListAllRegionsPartitionMapping(t *testing.T) {
	regions, err := ListAllRegions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test specific region-partition mappings
	testCases := []struct {
		regionId    string
		partitionId string
	}{
		{"us-east-1", "aws"},
		{"eu-west-1", "aws"},
		{"cn-north-1", "aws-cn"},
		{"us-gov-west-1", "aws-us-gov"},
	}

	regionPartitionMap := make(map[string]string)
	for _, region := range regions {
		regionPartitionMap[region.RegionId] = region.PartitionID
	}

	for _, tc := range testCases {
		if partition, found := regionPartitionMap[tc.regionId]; found {
			if partition != tc.partitionId {
				t.Errorf("Region %s expected to be in partition %s, but found in %s",
					tc.regionId, tc.partitionId, partition)
			}
		} else {
			t.Errorf("Region %s not found in results", tc.regionId)
		}
	}
}

func TestListAllRegionsServices(t *testing.T) {
	regions, err := ListAllRegions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Find us-east-1 and check it has services
	var usEast1 *RegionInfo
	for i, region := range regions {
		if region.RegionId == "us-east-1" {
			usEast1 = &regions[i]
			break
		}
	}

	if usEast1 == nil {
		t.Fatal("us-east-1 region not found")
	}

	if len(usEast1.Services) == 0 {
		t.Error("us-east-1 should have services available")
	}

	// Check that services are non-empty strings
	for i, service := range usEast1.Services {
		if service == "" {
			t.Errorf("Service at index %d in us-east-1 is empty string", i)
		}
	}

	// Check that services are sorted
	for i := 1; i < len(usEast1.Services); i++ {
		if usEast1.Services[i-1] > usEast1.Services[i] {
			t.Errorf("Services in us-east-1 are not sorted: %s > %s",
				usEast1.Services[i-1], usEast1.Services[i])
		}
	}
}

func TestListAllRegionsUniqueness(t *testing.T) {
	regions, err := ListAllRegions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that all region IDs are unique
	regionIds := make(map[string]bool)
	for _, region := range regions {
		if regionIds[region.RegionId] {
			t.Errorf("Duplicate region ID found: %s", region.RegionId)
		}
		regionIds[region.RegionId] = true
	}
}

func TestListAllRegionsStructure(t *testing.T) {
	regions, err := ListAllRegions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test that we have regions from multiple partitions
	partitionCounts := make(map[string]int)
	for _, region := range regions {
		partitionCounts[region.PartitionID]++
	}

	// AWS partition should have the most regions
	if partitionCounts["aws"] == 0 {
		t.Error("Expected regions from AWS partition")
	}

	// Should have more AWS regions than any other partition
	awsCount := partitionCounts["aws"]
	for partition, count := range partitionCounts {
		if partition != "aws" && count >= awsCount {
			t.Errorf("AWS partition should have more regions than %s: %d vs %d",
				partition, awsCount, count)
		}
	}

	// Check that we have a reasonable total number of regions
	if len(regions) < 10 {
		t.Errorf("Expected at least 10 regions, got %d", len(regions))
	}
}
