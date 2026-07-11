package main

import (
	"fmt"
	"log"

	"github.com/myerscode/aws-meta/pkg/services"
)

func main() {
	serviceMeta, err := services.ServiceMeta()
	if err != nil {
		log.Fatalf("Failed to get service metadata: %v", err)
	}

	fmt.Printf("Found %d services\n\n", len(serviceMeta))

	// Show some interesting services
	interestingServices := []string{"S3", "EC2", "Lambda", "RDS", "ACM"}
	for _, serviceId := range interestingServices {
		if service, found := serviceMeta[serviceId]; found {
			fmt.Printf("%s (%s)\n", service.ServiceFullName, service.ServiceId)
			fmt.Printf("  Protocol: %s, API Version: %s\n", service.Protocol, service.APIVersion)
			fmt.Printf("  Operations: %d\n\n", len(service.Operations))
		}
	}

	// Demonstrate ServiceMetaForRegion
	regionServices, err := services.ServiceMetaForRegion("eu-west-1")
	if err != nil {
		log.Fatalf("Failed to get services for eu-west-1: %v", err)
	}
	fmt.Printf("Services available in eu-west-1: %d\n", len(regionServices))
}
