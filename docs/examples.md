# Examples

This document provides usage examples for the AWS Metadata package functions.

## Services In Region

Get all services available in a specific AWS region:

```go
services, err := ServicesInRegion("us-east-1")
if err != nil {
    log.Fatal(err)
}

for _, service := range services {
    fmt.Println(service)
}
```

## Service Operations

Get all operations available for a specific AWS service by Service ID:

```go
operations, err := ServiceOperationsByServiceId("ACM")
if err != nil {
    log.Fatal(err)
}

for _, operation := range operations {
    fmt.Println(operation)
}
```

Or by Service Name:

```go
operations, err := ServiceOperationsByServiceName("AWS Certificate Manager")
if err != nil {
    log.Fatal(err)
}

for _, operation := range operations {
    fmt.Println(operation)
}
```

## Working Code Examples

See the `examples/` directory for complete, runnable code examples:

- `examples/list_all.go` - General usage examples
- `examples/services_in_region.go` - Demonstrates the `ServicesInRegion` function
- `examples/service_operations.go` - Demonstrates the `ServiceOperations` function
- `examples/partition_types.go` - Demonstrates partition categorization functions
- `examples/list_partitions.go` - Demonstrates detailed partition information
- `examples/list_regions.go` - Demonstrates comprehensive region information

Run any example with:
```bash
go run examples/example_name.go
```

## Partition Types

Get partitions categorized by type:

```go
import "github.com/myerscode/aws-meta/pkg/partitions"

// Get commercial partitions (standard AWS)
commercial, err := partitions.CommercialPartitions()
if err != nil {
    log.Fatal(err)
}

// Get sovereign partitions (China, GovCloud, EU Sovereign)
sovereign, err := partitions.SovereignPartitions()
if err != nil {
    log.Fatal(err)
}

// Get isolated partitions (air-gapped environments)
isolated, err := partitions.IsolatedPartitions()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Commercial: %v\n", commercial)
fmt.Printf("Sovereign: %v\n", sovereign)
fmt.Printf("Isolated: %v\n", isolated)
```

## Detailed Partition Information

Get comprehensive partition data including DNS suffixes and region details:

```go
import "github.com/myerscode/aws-meta/pkg/partitions"

partitionList, err := partitions.List()
if err != nil {
    log.Fatal(err)
}

for _, partition := range partitionList {
    fmt.Printf("Partition: %s\n", partition.ID)
    fmt.Printf("  DNS Suffix: %s\n", partition.DNSSuffix)
    fmt.Printf("  Region Pattern: %s\n", partition.RegionRegex)
    fmt.Printf("  Global Region: %s\n", partition.ImplicitGlobalRegion)
    fmt.Printf("  Regions: %d\n", len(partition.Regions))
    
    // Show first few regions
    for i, region := range partition.Regions {
        if i >= 3 { break } // Show only first 3
        fmt.Printf("    - %s (%s)\n", region.RegionId, region.RegionName)
    }
}
```

## Comprehensive Region Information

Get detailed information about all regions including partition mapping and service counts:

```go
import "github.com/myerscode/aws-meta/pkg/regions"

regionList, err := regions.ListAllRegions()
if err != nil {
    log.Fatal(err)
}

for _, region := range regionList {
    fmt.Printf("Region: %s (%s)\n", region.RegionId, region.RegionName)
    fmt.Printf("  Partition: %s\n", region.PartitionID)
    fmt.Printf("  Services: %d\n", len(region.Services))
    
    // Show first few services
    for i, service := range region.Services {
        if i >= 3 { break } // Show only first 3
        fmt.Printf("    - %s\n", service)
    }
}
```

## Complete Example

Here's a complete example showing how to use multiple functions together:

```go
package main

import (
    "fmt"
    "log"

    "github.com/myerscode/aws-meta/pkg/partitions"
    "github.com/myerscode/aws-meta/pkg/services"
)

func main() {
    // Get all partitions
    partitionNames := partitions.AllPartitionNames()
    fmt.Printf("Found %d partitions\n", len(partitionNames))

    // Get all regions
    regions := services.AllRegionNames()
    fmt.Printf("Found %d regions\n", len(regions))

    // Get services in a specific region
    servicesInRegion, err := services.ServicesInRegion("us-east-1")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d services in us-east-1\n", len(servicesInRegion))

    // Get operations for a specific service
    operations, err := services.ServiceOperations("S3")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d operations for S3\n", len(operations))
    
    // Show first few operations
    for i, op := range operations {
        if i >= 5 { // Show only first 5
            break
        }
        fmt.Printf("  - %s\n", op)
    }
}
```