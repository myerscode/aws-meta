# API Reference

This document provides an overview of all available packages and their methods in the AWS Metadata library.

## Package Documentation

The AWS Metadata library is organized into focused packages, each providing specific functionality:

### [Partitions Package](partitions.md)
**Import:** `github.com/myerscode/aws-meta/pkg/partitions`

Functions for working with AWS partitions:
- `AllPartitionNames()` - Get all partition names
- `List()` - Get detailed partition information
- `CommercialPartitions()` - Get commercial partitions
- `SovereignPartitions()` - Get sovereign partitions  
- `IsolatedPartitions()` - Get isolated partitions

### [Regions Package](regions.md)
**Import:** `github.com/myerscode/aws-meta/pkg/regions`

Functions for working with AWS regions:
- `ListAllRegions()` - Get comprehensive region information with partition mapping and service availability

### [Services Package](services.md)
**Import:** `github.com/myerscode/aws-meta/pkg/services`

Functions for working with AWS services and regions:
- `ServiceMeta()` - Get detailed service metadata for all services
- `ServiceMetaForRegion()` - Get detailed service metadata for services in a specific region
- `AllRegionNames()` - Get all region names
- `AllServiceNames()` - Get all service names
- `ServicesInRegion()` - Get services available in a region
- `ServiceOperationsByServiceId()` - Get operations by service ID
- `ServiceOperationsByServiceName()` - Get operations by service name
- `ServiceOperations()` - Alias for ServiceOperationsByServiceId

## Quick Start Examples

### Working with Partitions
```go
import "github.com/myerscode/aws-meta/pkg/partitions"

// Get all partition names
partitionNames := partitions.AllPartitionNames()
fmt.Printf("Available partitions: %v\n", partitionNames)

// Get detailed partition information
partitionList, err := partitions.List()
if err != nil {
    log.Fatal(err)
}
```

### Working with Regions
```go
import "github.com/myerscode/aws-meta/pkg/regions"

// Get comprehensive region information
regionList, err := regions.ListAllRegions()
if err != nil {
    log.Fatal(err)
}

for _, region := range regionList {
    fmt.Printf("Region: %s (%s)\n", region.RegionId, region.RegionName)
    fmt.Printf("  Partition: %s\n", region.PartitionID)
    fmt.Printf("  Services: %d\n", len(region.Services))
}
```

### Working with Services and Regions
```go
import "github.com/myerscode/aws-meta/pkg/services"

// Get detailed service metadata
allServices, err := services.ServiceMeta()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total services: %d\n", len(allServices))

// Get detailed service metadata for a specific region
regionServices, err := services.ServiceMetaForRegion("us-east-1")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Services available in us-east-1: %d\n", len(regionServices))

// Get all regions
regions := services.AllRegionNames()
fmt.Printf("Total regions: %d\n", len(regions))
for _, region := range regions[:5] { // Show first 5
    fmt.Printf("  - %s\n", region)
}
```

For detailed documentation on each package's functions, see the individual package documentation pages linked above.

## Data Sources

All data is sourced from the official AWS Botocore library, ensuring accuracy and up-to-date information:

- **Partitions**: From `botocore/data/partitions.json`
- **Regions**: From `botocore/data/endpoints.json`
- **Services**: From `botocore/data/*/service-*.json` files
- **Operations**: Extracted from individual service definition files

## Error Handling

All functions that return errors follow Go conventions:

- **Success**: `error` is `nil`
- **Failure**: `error` contains a descriptive message
- **Input validation**: Invalid inputs return specific error messages
- **Case sensitivity**: All inputs are case-sensitive and must match AWS conventions

## Performance Notes

- All data is loaded from pre-generated manifest files
- Functions return quickly as data is already parsed and sorted
- No network calls are made during function execution
- Data is cached in memory after first access

## Complete Examples

For complete, runnable examples, see the `examples/` directory:
- `examples/list-all/` - General usage examples
- `examples/services-in-region/` - Regional service availability
- `examples/service-operations/` - Service operation listings
- `examples/partition-types/` - Partition categorization
- `examples/list-partitions/` - Detailed partition information
- `examples/list-regions/` - Comprehensive region information