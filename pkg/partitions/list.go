package partitions

import (
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

// ListPartitions returns detailed information about all AWS partitions.
// This includes partition metadata such as DNS suffixes, region regex patterns,
// and the list of regions within each partition.
//
// Returns:
//   - aws.PartitionSchemas: A slice of partition schemas containing detailed partition information
//   - error: An error if there's an issue reading the manifest
//
// Example:
//
//	partitions, err := partitions.ListPartitions()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, partition := range partitions {
//	    fmt.Printf("Partition: %s\n", partition.ID)
//	    fmt.Printf("  DNS Suffix: %s\n", partition.DNSSuffix)
//	    fmt.Printf("  Regions: %d\n", len(partition.Regions))
//	}
func ListPartitions() (aws.PartitionSchemas, error) {
	return partitionManifest()
}
