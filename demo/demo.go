package main

import (
	"fmt"
	"github.com/myerscode/aws-meta/internal/util"
	"github.com/myerscode/aws-meta/pkg/services"
)

func main() {

	serviceNames := services.AllServiceNames()

	util.LogInfo(fmt.Sprintf("AllServiceNames() returned %d service names", len(serviceNames)))

	regionNames := services.AllRegionNames()

	util.LogInfo(fmt.Sprintf("AllRegionNames() returned %d region names", len(regionNames)))

	partitionNames := services.AllPartitionNames()

	util.LogInfo(fmt.Sprintf("AllPartitionNames() returned %d partition names", len(partitionNames)))
}
