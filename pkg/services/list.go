package services

import (
	"github.com/myerscode/aws-meta/internal/aws"
	"github.com/myerscode/aws-meta/pkg/data"
)

func partitionManifest() (aws.PartitionSchemas, error) {
	manifest, err := data.PartitionManifest()

	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func regionManifest() (aws.RegionSchemas, error) {
	manifest, err := data.RegionsManifest()

	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func serviceManifest() (aws.ServiceSchemas, error) {
	manifest, err := data.ServiceManifest()

	if err != nil {
		return nil, err
	}

	return manifest, nil
}

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

func AllRegionNames() []string {
	manifest, err := regionManifest()

	var regionNames []string

	if err != nil {
		return regionNames
	}

	for _, meta := range manifest {
		for _, region := range meta.Regions {
			regionNames = append(regionNames, region.RegionName)
		}
	}

	return regionNames
}

func AllServiceNames() []string {

	manifest, err := serviceManifest()

	var serviceNames []string

	if err != nil {
		return serviceNames
	}

	for _, meta := range manifest {
		serviceNames = append(serviceNames, meta.ServiceFullName)
	}

	return serviceNames
}
