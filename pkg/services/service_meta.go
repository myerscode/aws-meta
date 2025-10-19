package services

import (
	"github.com/myerscode/aws-meta/internal/aws"
)

// ServiceMeta returns detailed information about all AWS services.
// This includes service metadata such as service IDs, full names, operations,
// and regional availability information.
//
// Returns:
//   - aws.ServiceSchemas: A slice of service schemas containing detailed service information
//   - error: An error if there's an issue reading the manifest
//
// Example:
//
//	services, err := services.ServiceMeta()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, service := range services {
//	    fmt.Printf("Service: %s\n", service.ServiceFullName)
//	    fmt.Printf("  ID: %s\n", service.ServiceId)
//	    fmt.Printf("  Operations: %d\n", len(service.Operations))
//	}
func ServiceMeta() (aws.ServiceSchemas, error) {
	return serviceManifest()
}

// ServiceMetaForRegion returns detailed information about AWS services available in the specified region.
// This includes service metadata such as service IDs, full names, operations, and protocols
// for only the services that are available in the given region.
//
// Parameters:
//   - regionName: The AWS region name (e.g., "us-east-1", "eu-west-1")
//
// Returns:
//   - aws.ServiceSchemas: A map of service schemas containing detailed service information for services available in the region
//   - error: An error if the region name is invalid or if there's an issue reading the manifest
//
// Example:
//
//	services, err := services.ServiceMetaForRegion("us-east-1")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for serviceId, service := range services {
//	    fmt.Printf("Service: %s (available in us-east-1)\n", service.ServiceFullName)
//	    fmt.Printf("  ID: %s\n", service.ServiceId)
//	    fmt.Printf("  Operations: %d\n", len(service.Operations))
//	}
func ServiceMetaForRegion(regionName string) (aws.ServiceSchemas, error) {
	// First get the list of service endpoint prefixes available in the region
	serviceEndpointsInRegion, err := ServicesInRegion(regionName)
	if err != nil {
		return nil, err
	}

	// Get all service metadata
	allServiceMeta, err := ServiceMeta()
	if err != nil {
		return nil, err
	}

	// Create a map to quickly lookup service endpoint prefixes
	endpointPrefixSet := make(map[string]bool)
	for _, endpointPrefix := range serviceEndpointsInRegion {
		endpointPrefixSet[endpointPrefix] = true
	}

	// Filter service metadata to only include services available in the region
	// Match by EndpointPrefix field
	regionServiceMeta := make(aws.ServiceSchemas)
	for serviceId, service := range allServiceMeta {
		if endpointPrefixSet[service.EndpointPrefix] {
			regionServiceMeta[serviceId] = service
		}
	}

	return regionServiceMeta, nil
}
