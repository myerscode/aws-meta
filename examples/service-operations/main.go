package main

import (
	"fmt"
	"log"

	"github.com/myerscode/aws-meta/examples/shared"
	"github.com/myerscode/aws-meta/pkg/services"
)

func main() {
	testServiceIds := []string{"ACM", "S3", "EC2", "Lambda"}

	for _, serviceId := range testServiceIds {
		operations, err := services.ServiceOperationsByServiceId(serviceId)
		if err != nil {
			log.Printf("Failed to get operations for %s: %v", serviceId, err)
			continue
		}
		fmt.Printf("%s: %d operations", serviceId, len(operations))
		if len(operations) > 0 {
			fmt.Printf(" (first 5: %v)", operations[:shared.Min(5, len(operations))])
		}
		fmt.Println()
	}

	// Test error case
	_, err := services.ServiceOperationsByServiceId("INVALID-SERVICE")
	if err != nil {
		fmt.Printf("\nExpected error for invalid service: %v\n", err)
	}
}
