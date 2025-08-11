package services

import (
	"strings"
)

// IsolatedPartitions returns a list of AWS isolated partition names.
// Isolated partitions are air-gapped environments for highly sensitive workloads.
// These typically include partitions with "iso" in their name.
//
// Returns:
//   - []string: A sorted list of isolated partition names (e.g., "aws-iso", "aws-iso-b", "aws-iso-e", "aws-iso-f")
//
// Example:
//
//	partitions, err := IsolatedPartitions()
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
//	partitions, err := SovereignPartitions()
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
//	partitions, err := CommercialPartitions()
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
