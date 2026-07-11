package main

import (
	"fmt"
	"log"

	"github.com/myerscode/aws-meta/pkg/partitions"
)

func main() {
	commercial, err := partitions.CommercialPartitions()
	if err != nil {
		log.Fatalf("Failed to get commercial partitions: %v", err)
	}
	fmt.Printf("Commercial partitions (%d): %v\n", len(commercial), commercial)

	sovereign, err := partitions.SovereignPartitions()
	if err != nil {
		log.Fatalf("Failed to get sovereign partitions: %v", err)
	}
	fmt.Printf("Sovereign partitions (%d): %v\n", len(sovereign), sovereign)

	isolated, err := partitions.IsolatedPartitions()
	if err != nil {
		log.Fatalf("Failed to get isolated partitions: %v", err)
	}
	fmt.Printf("Isolated partitions (%d): %v\n", len(isolated), isolated)

	allPartitions := partitions.AllPartitionNames()
	fmt.Printf("\nTotal: %d partitions\n", len(allPartitions))
}
