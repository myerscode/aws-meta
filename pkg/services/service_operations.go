package services

import (
	"fmt"
)

// ServiceOperationsByServiceId returns a list of operation names available for the specified service ID.
// The returned operations are sorted alphabetically.
//
// Parameters:
//   - serviceId: The AWS service ID (e.g., "ACM", "S3", "EC2")
//
// Returns:
//   - []string: A sorted list of operation names available for the service
//   - error: An error if the service ID is invalid or if there's an issue reading the manifest
//
// Example:
//
//	operations, err := ServiceOperationsByServiceId("ACM")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d operations for ACM\n", len(operations))
func ServiceOperationsByServiceId(serviceId string) ([]string, error) {
	manifest, err := serviceManifest()
	if err != nil {
		return nil, err
	}

	// Look for the service by ServiceId
	for _, service := range manifest {
		if service.ServiceId == serviceId {
			return service.Operations, nil
		}
	}

	// If we get here, the service was not found
	return nil, fmt.Errorf("invalid service ID: %s", serviceId)
}

// ServiceOperationsByServiceName returns a list of operation names available for the specified service name.
// The returned operations are sorted alphabetically.
//
// Parameters:
//   - serviceName: The AWS service full name (e.g., "AWS Certificate Manager", "Amazon Simple Storage Service", "Amazon Elastic Compute Cloud")
//
// Returns:
//   - []string: A sorted list of operation names available for the service
//   - error: An error if the service name is invalid or if there's an issue reading the manifest
//
// Example:
//
//	operations, err := ServiceOperationsByServiceName("AWS Certificate Manager")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d operations for AWS Certificate Manager\n", len(operations))
func ServiceOperationsByServiceName(serviceName string) ([]string, error) {
	manifest, err := serviceManifest()
	if err != nil {
		return nil, err
	}

	// Look for the service by ServiceFullName
	for _, service := range manifest {
		if service.ServiceFullName == serviceName {
			return service.Operations, nil
		}
	}

	// If we get here, the service was not found
	return nil, fmt.Errorf("invalid service name: %s", serviceName)
}

// ServiceOperations returns a list of operation names available for the specified service ID.
// This is an alias for ServiceOperationsByServiceId for convenience.
// The returned operations are sorted alphabetically.
//
// Parameters:
//   - serviceId: The AWS service ID (e.g., "ACM", "S3", "EC2")
//
// Returns:
//   - []string: A sorted list of operation names available for the service
//   - error: An error if the service ID is invalid or if there's an issue reading the manifest
//
// Example:
//
//	operations, err := ServiceOperations("ACM")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d operations for ACM\n", len(operations))
func ServiceOperations(serviceId string) ([]string, error) {
	return ServiceOperationsByServiceId(serviceId)
}
