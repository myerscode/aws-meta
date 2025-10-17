package partitions

import (
	"testing"

	"github.com/myerscode/aws-meta/internal/aws"
)

func TestListPartitions(t *testing.T) {
	partitions, err := ListPartitions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(partitions) == 0 {
		t.Fatal("Expected at least one partition")
	}

	// Check that all partitions have required fields
	for i, partition := range partitions {
		if partition.ID == "" {
			t.Errorf("Partition at index %d has empty ID", i)
		}
		if partition.DNSSuffix == "" {
			t.Errorf("Partition %s has empty DNSSuffix", partition.ID)
		}
		if partition.RegionRegex == "" {
			t.Errorf("Partition %s has empty RegionRegex", partition.ID)
		}
		if len(partition.Regions) == 0 {
			t.Errorf("Partition %s has no regions", partition.ID)
		}

		// Check that all regions have required fields
		for j, region := range partition.Regions {
			if region.RegionId == "" {
				t.Errorf("Partition %s, region at index %d has empty RegionId", partition.ID, j)
			}
			if region.RegionName == "" {
				t.Errorf("Partition %s, region %s has empty RegionName", partition.ID, region.RegionId)
			}
		}
	}
}

func TestListPartitionsContent(t *testing.T) {
	partitions, err := ListPartitions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check for known partitions
	expectedPartitions := []string{"aws", "aws-cn", "aws-us-gov", "aws-eusc", "aws-iso", "aws-iso-b", "aws-iso-e", "aws-iso-f"}
	foundPartitions := make(map[string]bool)

	for _, partition := range partitions {
		foundPartitions[partition.ID] = true
	}

	for _, expected := range expectedPartitions {
		if !foundPartitions[expected] {
			t.Errorf("Expected partition %s not found", expected)
		}
	}

	// Check that partitions are sorted by ID
	for i := 1; i < len(partitions); i++ {
		if partitions[i-1].ID > partitions[i].ID {
			t.Errorf("Partitions are not sorted: %s > %s", partitions[i-1].ID, partitions[i].ID)
		}
	}
}

func TestListPartitionsSpecificData(t *testing.T) {
	partitions, err := ListPartitions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Find and test specific partition data
	var awsPartition *aws.PartitionSchema
	var awsCnPartition *aws.PartitionSchema

	for i, partition := range partitions {
		if partition.ID == "aws" {
			awsPartition = &partitions[i]
		}
		if partition.ID == "aws-cn" {
			awsCnPartition = &partitions[i]
		}
	}

	// Test AWS partition
	if awsPartition == nil {
		t.Fatal("AWS partition not found")
	}
	if awsPartition.DNSSuffix != "amazonaws.com" {
		t.Errorf("Expected AWS partition DNS suffix to be 'amazonaws.com', got %s", awsPartition.DNSSuffix)
	}
	if awsPartition.ImplicitGlobalRegion != "us-east-1" {
		t.Errorf("Expected AWS partition implicit global region to be 'us-east-1', got %s", awsPartition.ImplicitGlobalRegion)
	}

	// Test AWS China partition
	if awsCnPartition == nil {
		t.Fatal("AWS China partition not found")
	}
	if awsCnPartition.DNSSuffix != "amazonaws.com.cn" {
		t.Errorf("Expected AWS China partition DNS suffix to be 'amazonaws.com.cn', got %s", awsCnPartition.DNSSuffix)
	}

	// Test that AWS partition has more regions than China partition
	if len(awsPartition.Regions) <= len(awsCnPartition.Regions) {
		t.Errorf("Expected AWS partition to have more regions than China partition: %d vs %d",
			len(awsPartition.Regions), len(awsCnPartition.Regions))
	}
}

func TestListPartitionsRegionData(t *testing.T) {
	partitions, err := ListPartitions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Find AWS partition and test some known regions
	var awsPartition *aws.PartitionSchema
	for i, partition := range partitions {
		if partition.ID == "aws" {
			awsPartition = &partitions[i]
			break
		}
	}

	if awsPartition == nil {
		t.Fatal("AWS partition not found")
	}

	// Check for some known regions
	expectedRegions := []string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-1"}
	foundRegions := make(map[string]bool)

	for _, region := range awsPartition.Regions {
		foundRegions[region.RegionId] = true
	}

	for _, expected := range expectedRegions {
		if !foundRegions[expected] {
			t.Errorf("Expected region %s not found in AWS partition", expected)
		}
	}

	// Check that regions within a partition are sorted
	for i := 1; i < len(awsPartition.Regions); i++ {
		if awsPartition.Regions[i-1].RegionId > awsPartition.Regions[i].RegionId {
			t.Errorf("Regions in AWS partition are not sorted: %s > %s",
				awsPartition.Regions[i-1].RegionId, awsPartition.Regions[i].RegionId)
		}
	}
}
