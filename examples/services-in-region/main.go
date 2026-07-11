package main

import (
	"fmt"
	"log"

	"github.com/myerscode/aws-meta/examples/shared"
	"github.com/myerscode/aws-meta/pkg/services"
)

func main() {
	servicesUSEast1, err := services.ServicesInRegion("us-east-1")
	if err != nil {
		log.Fatalf("Failed to get services for us-east-1: %v", err)
	}
	fmt.Printf("us-east-1: %d services\n", len(servicesUSEast1))
	if len(servicesUSEast1) > 0 {
		fmt.Printf("  First 5: %v\n", servicesUSEast1[:shared.Min(5, len(servicesUSEast1))])
	}

	servicesEUWest1, err := services.ServicesInRegion("eu-west-1")
	if err != nil {
		log.Fatalf("Failed to get services for eu-west-1: %v", err)
	}
	fmt.Printf("eu-west-1: %d services\n", len(servicesEUWest1))

	// Test invalid region
	_, err = services.ServicesInRegion("invalid-region")
	if err != nil {
		fmt.Printf("\nExpected error for invalid region: %v\n", err)
	}
}
