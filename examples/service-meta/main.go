package main

import (
	"fmt"

	"github.com/myerscode/aws-meta/internal/util"
	"github.com/myerscode/aws-meta/pkg/services"
)

func main() {
	util.LogInfo("=== AWS Service Metadata ===")

	serviceMeta, err := services.ServiceMeta()
	if err != nil {
		util.LogError(fmt.Sprintf("Failed to get service metadata: %v", err))
		return
	}

	util.LogInfo(fmt.Sprintf("Found %d services:", len(serviceMeta)))
	fmt.Println()

	// Show some interesting services with their metadata
	interestingServices := []string{"S3", "EC2", "Lambda", "RDS", "ACM"}

	for _, serviceId := range interestingServices {
		if service, found := serviceMeta[serviceId]; found {
			util.LogInfo(fmt.Sprintf("Service: %s", service.ServiceFullName))
			util.LogInfo(fmt.Sprintf("  ID: %s", service.ServiceId))
			util.LogInfo(fmt.Sprintf("  Protocol: %s", service.Protocol))
			util.LogInfo(fmt.Sprintf("  API Version: %s", service.APIVersion))
			util.LogInfo(fmt.Sprintf("  Endpoint Prefix: %s", service.EndpointPrefix))
			util.LogInfo(fmt.Sprintf("  Operations: %d", len(service.Operations)))

			// Show first few operations
			if len(service.Operations) > 0 {
				util.LogInfo("  Sample operations:")
				maxOps := 5
				if len(service.Operations) < maxOps {
					maxOps = len(service.Operations)
				}
				for i := 0; i < maxOps; i++ {
					util.LogInfo(fmt.Sprintf("    - %s", service.Operations[i]))
				}
				if len(service.Operations) > maxOps {
					util.LogInfo(fmt.Sprintf("    ... and %d more", len(service.Operations)-maxOps))
				}
			}
			fmt.Println()
		}
	}

	// Demonstrate ServiceMetaForRegion
	util.LogInfo("=== Services Available in eu-west-1 ===")
	regionServices, err := services.ServiceMetaForRegion("eu-west-1")
	if err != nil {
		util.LogError(fmt.Sprintf("Failed to get services for eu-west-1: %v", err))
		return
	}

	util.LogInfo(fmt.Sprintf("Found %d services available in eu-west-1", len(regionServices)))

	// Show a few services available in eu-west-1
	count := 0
	for serviceId, service := range regionServices {
		if count >= 5 { // Show only first 5
			break
		}
		util.LogInfo(fmt.Sprintf("  %s: %s", serviceId, service.ServiceFullName))
		count++
	}
	if len(regionServices) > 5 {
		util.LogInfo(fmt.Sprintf("  ... and %d more services", len(regionServices)-5))
	}
	fmt.Println()

	// Summary statistics
	util.LogInfo("=== Service Statistics ===")
	totalOperations := 0
	protocolCounts := make(map[string]int)

	for _, service := range serviceMeta {
		totalOperations += len(service.Operations)
		if service.Protocol != "" {
			protocolCounts[service.Protocol]++
		}
	}

	util.LogInfo(fmt.Sprintf("Total services: %d", len(serviceMeta)))
	util.LogInfo(fmt.Sprintf("Total operations: %d", totalOperations))
	util.LogInfo(fmt.Sprintf("Average operations per service: %.1f", float64(totalOperations)/float64(len(serviceMeta))))

	util.LogInfo("Protocol distribution:")
	for protocol, count := range protocolCounts {
		util.LogInfo(fmt.Sprintf("  %s: %d services", protocol, count))
	}

	// Find service with most operations
	var maxOpsService string
	var maxOps int
	for serviceId, service := range serviceMeta {
		if len(service.Operations) > maxOps {
			maxOps = len(service.Operations)
			maxOpsService = serviceId
		}
	}

	if maxOpsService != "" {
		service := serviceMeta[maxOpsService]
		util.LogInfo(fmt.Sprintf("Service with most operations: %s (%s) - %d operations",
			maxOpsService, service.ServiceFullName, maxOps))
	}
}
