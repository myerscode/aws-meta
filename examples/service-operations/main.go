package main

import (
	"fmt"

	"github.com/myerscode/aws-meta/examples/shared"
	"github.com/myerscode/aws-meta/internal/util"
	"github.com/myerscode/aws-meta/pkg/services"
)

func main() {
	// Test with valid services by Service ID
	util.LogInfo("=== Testing ServiceOperationsByServiceId ===")
	testServiceIds := []string{"ACM", "S3", "EC2", "Lambda"}

	for _, serviceId := range testServiceIds {
		util.LogInfo(fmt.Sprintf("Testing ServiceOperationsByServiceId with %s:", serviceId))
		operations, err := services.ServiceOperationsByServiceId(serviceId)
		if err != nil {
			util.LogError(fmt.Sprintf("Failed to get operations for %s: %v", serviceId, err))
			continue
		}
		util.LogInfo(fmt.Sprintf("Found %d operations for %s", len(operations), serviceId))
		if len(operations) > 0 {
			util.LogInfo(fmt.Sprintf("First few operations: %v", operations[:shared.Min(5, len(operations))]))
		}
		fmt.Println() // Add spacing between services
	}

	// Test with valid services by Service Name
	util.LogInfo("=== Testing ServiceOperationsByServiceName ===")
	testServiceNames := []string{
		"AWS Certificate Manager",
		"Amazon Simple Storage Service",
		"Amazon Elastic Compute Cloud",
		"AWS Lambda",
	}

	for _, serviceName := range testServiceNames {
		util.LogInfo(fmt.Sprintf("Testing ServiceOperationsByServiceName with '%s':", serviceName))
		operations, err := services.ServiceOperationsByServiceName(serviceName)
		if err != nil {
			util.LogError(fmt.Sprintf("Failed to get operations for '%s': %v", serviceName, err))
			continue
		}
		util.LogInfo(fmt.Sprintf("Found %d operations for '%s'", len(operations), serviceName))
		if len(operations) > 0 {
			util.LogInfo(fmt.Sprintf("First few operations: %v", operations[:shared.Min(3, len(operations))]))
		}
		fmt.Println() // Add spacing between services
	}

	// Test error cases
	util.LogInfo("=== Testing Error Cases ===")

	// Test with invalid service ID
	util.LogInfo("Testing with invalid service ID:")
	_, err := services.ServiceOperationsByServiceId("INVALID-SERVICE")
	if err != nil {
		util.LogInfo(fmt.Sprintf("Expected error: %v", err))
	} else {
		util.LogError("ERROR: Should have gotten an error for invalid service ID")
	}

	// Test with invalid service name
	util.LogInfo("Testing with invalid service name:")
	_, err = services.ServiceOperationsByServiceName("Invalid Service Name")
	if err != nil {
		util.LogInfo(fmt.Sprintf("Expected error: %v", err))
	} else {
		util.LogError("ERROR: Should have gotten an error for invalid service name")
	}

	// Test case sensitivity for service ID
	util.LogInfo("Testing case sensitivity for service ID:")
	_, err = services.ServiceOperationsByServiceId("acm")
	if err != nil {
		util.LogInfo(fmt.Sprintf("Expected error (case sensitive): %v", err))
	} else {
		util.LogError("ERROR: Should have gotten an error for lowercase service ID")
	}

	// Test case sensitivity for service name
	util.LogInfo("Testing case sensitivity for service name:")
	_, err = services.ServiceOperationsByServiceName("aws certificate manager")
	if err != nil {
		util.LogInfo(fmt.Sprintf("Expected error (case sensitive): %v", err))
	} else {
		util.LogError("ERROR: Should have gotten an error for lowercase service name")
	}

	// Test ServiceOperations alias function
	util.LogInfo("=== Testing ServiceOperations Alias ===")
	util.LogInfo("Testing ServiceOperations function (alias for ServiceOperationsByServiceId):")
	operations, err := services.ServiceOperations("ACM")
	if err != nil {
		util.LogError(fmt.Sprintf("Failed to get operations using ServiceOperations: %v", err))
	} else {
		util.LogInfo(fmt.Sprintf("ServiceOperations alias works: Found %d operations for ACM", len(operations)))
	}
}
