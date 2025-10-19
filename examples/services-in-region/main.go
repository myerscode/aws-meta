package main

import (
	"fmt"
	"os"

	"github.com/myerscode/aws-meta/examples/shared"
	"github.com/myerscode/aws-meta/internal/util"
	"github.com/myerscode/aws-meta/pkg/services"
)

func main() {
	// Test with a valid region
	util.LogInfo("Testing ServicesInRegion with us-east-1:")
	servicesUSEast1, err := services.ServicesInRegion("us-east-1")
	if err != nil {
		util.LogError(fmt.Sprintf("Failed to get services for us-east-1: %v", err))
		os.Exit(1)
	}
	util.LogInfo(fmt.Sprintf("Found %d services in us-east-1", len(servicesUSEast1)))
	if len(servicesUSEast1) > 0 {
		util.LogInfo(fmt.Sprintf("First few services: %v", servicesUSEast1[:shared.Min(5, len(servicesUSEast1))]))
	}

	util.LogInfo("Testing ServicesInRegion with eu-west-1:")
	servicesEUWest1, err := services.ServicesInRegion("eu-west-1")
	if err != nil {
		util.LogError(fmt.Sprintf("Failed to get services for eu-west-1: %v", err))
		os.Exit(1)
	}
	util.LogInfo(fmt.Sprintf("Found %d services in eu-west-1", len(servicesEUWest1)))

	// Test with invalid region
	util.LogInfo("Testing ServicesInRegion with invalid region:")
	_, err = services.ServicesInRegion("invalid-region")
	if err != nil {
		util.LogInfo(fmt.Sprintf("Expected error: %v", err))
	} else {
		util.LogError("ERROR: Should have gotten an error for invalid region")
	}
}
