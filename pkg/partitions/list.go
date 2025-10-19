package partitions

import (
	"strings"

	"github.com/myerscode/aws-meta/internal/aws"
	"github.com/myerscode/aws-meta/pkg/data"
)

// partitionManifest loads the partition manifest data
func partitionManifest() (aws.PartitionSchemas, error) {
	manifest, err := data.PartitionManifest()
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

// List returns detailed information about all AWS partitions.
// This includes partition metadata such as DNS suffixes, region regex patterns,
// and the list of regions within each partition.
//
// Returns:
//   - aws.PartitionSchemas: A slice of partition schemas containing detailed partition information
//   - error: An error if there's an issue reading the manifest
//
// Example:
//
//	partitions, err := partitions.List()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, partition := range partitions {
//	    fmt.Printf("Partition: %s\n", partition.ID)
//	    fmt.Printf("  DNS Suffix: %s\n", partition.DNSSuffix)
//	    fmt.Printf("  Regions: %d\n", len(partition.Regions))
//	}
func List() (aws.PartitionSchemas, error) {
	return partitionManifest()
}

// AllPartitionNames returns a list of all AWS partition names.
//
// Returns:
//   - []string: A list of all partition names (e.g., "aws", "aws-cn", "aws-us-gov")
//
// Example:
//
//	partitions := partitions.AllPartitionNames()
//	fmt.Printf("Available partitions: %v\n", partitions)
func AllPartitionNames() []string {
	manifest, err := partitionManifest()

	var partitionNames []string

	if err != nil {
		return partitionNames
	}

	for _, meta := range manifest {
		partitionNames = append(partitionNames, meta.ID)
	}

	return partitionNames
}

// IsolatedPartitions returns a list of AWS isolated partition names.
// Isolated partitions are air-gapped environments for highly sensitive workloads.
// These typically include partitions with "iso" in their name.
//
// Returns:
//   - []string: A sorted list of isolated partition names (e.g., "aws-iso", "aws-iso-b", "aws-iso-e", "aws-iso-f")
//
// Example:
//
//	partitions, err := partitions.IsolatedPartitions()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d isolated partitions\n", len(partitions))
func IsolatedPartitions() ([]string, error) {
	manifest, err := partitionManifest()
	if err != nil {
		return nil, err
	}

	var isolatedPartitions []string
	for _, partition := range manifest {
		// Isolated partitions contain "iso" in their ID
		if strings.Contains(partition.ID, "iso") {
			isolatedPartitions = append(isolatedPartitions, partition.ID)
		}
	}

	return isolatedPartitions, nil
}

// SovereignPartitions returns a list of AWS sovereign partition names.
// Sovereign partitions are operated by local entities to meet specific regulatory requirements.
// These include China, GovCloud, and EU Sovereign Cloud partitions.
//
// Returns:
//   - []string: A sorted list of sovereign partition names (e.g., "aws-cn", "aws-us-gov", "aws-eusc")
//
// Example:
//
//	partitions, err := partitions.SovereignPartitions()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d sovereign partitions\n", len(partitions))
func SovereignPartitions() ([]string, error) {
	manifest, err := partitionManifest()
	if err != nil {
		return nil, err
	}

	var sovereignPartitions []string
	for _, partition := range manifest {
		// Sovereign partitions are non-commercial, non-isolated partitions
		if partition.ID != "aws" && !strings.Contains(partition.ID, "iso") {
			sovereignPartitions = append(sovereignPartitions, partition.ID)
		}
	}

	return sovereignPartitions, nil
}

// CommercialPartitions returns a list of AWS commercial partition names.
// Commercial partitions are the standard AWS partitions available to the general public.
// This typically includes just the main "aws" partition.
//
// Returns:
//   - []string: A sorted list of commercial partition names (e.g., "aws")
//
// Example:
//
//	partitions, err := partitions.CommercialPartitions()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d commercial partitions\n", len(partitions))
func CommercialPartitions() ([]string, error) {
	manifest, err := partitionManifest()
	if err != nil {
		return nil, err
	}

	var commercialPartitions []string
	for _, partition := range manifest {
		// Commercial partition is the standard "aws" partition
		if partition.ID == "aws" {
			commercialPartitions = append(commercialPartitions, partition.ID)
		}
	}

	return commercialPartitions, nil
}
