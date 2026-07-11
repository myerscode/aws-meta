package main

import (
	"fmt"
	"log"

	"github.com/myerscode/aws-meta/pkg/partitions"
	"github.com/myerscode/aws-meta/pkg/services"
)

func main() {
	serviceNames := services.AllServiceNames()
	log.Printf("AllServiceNames() returned %d service names", len(serviceNames))

	regionNames := services.AllRegionNames()
	log.Printf("AllRegionNames() returned %d region names", len(regionNames))

	partitionNames := partitions.AllPartitionNames()
	fmt.Printf("AllPartitionNames() returned %d partition names\n", len(partitionNames))
}
