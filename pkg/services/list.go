package services

import (
	"github.com/myerscode/aws-meta/internal/aws"
	"github.com/myerscode/aws-meta/pkg/data"
)

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

// List returns detailed information about all AWS services.
// This includes service metadata such as service IDs, full names, operations,
// and regional availability information.
//
// Returns:
//   - aws.ServiceSchemas: A map of service schemas containing detailed service information
//   - error: An error if there's an issue reading the manifest
//
// Example:
//
//	services, err := services.List()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for serviceId, service := range services {
//	    fmt.Printf("Service: %s\n", service.ServiceFullName)
//	    fmt.Printf("  ID: %s\n", service.ServiceId)
//	    fmt.Printf("  Operations: %d\n", len(service.Operations))
//	}
func List() (aws.ServiceSchemas, error) {
	return serviceManifest()
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
