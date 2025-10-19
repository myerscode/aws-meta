# AWS Metadata
> A package for exposing information about AWS Partitions, Regions and Services


Knowing what AWS services are available in which regions and partitions can be a bit of a pain; 
this app aims to keep an up to date reference for the available endpoints for all AWS services.

## Features

- **Partitions**: Get all AWS partition names (aws, aws-cn, aws-us-gov)
- **Regions**: List all AWS regions across all partitions
- **Services**: Access AWS service information and operations
- **Regional Services**: Find which services are available in specific regions
- **Service Operations**: Get all available operations for any AWS service

## Quick Start

```go
import "github.com/myerscode/aws-meta/pkg/services"

// Get all regions
regions := services.AllRegionNames()

// Get services in a region
servicesInRegion, err := services.ServicesInRegion("us-east-1")

// Get operations for a service
operations, err := services.ServiceOperations("S3")
```

## Documentation

- [API Reference](docs/methods.md) - Overview of all packages and functions
- [Partitions Package](docs/partitions.md) - Partition-related functions
- [Services Package](docs/services.md) - Service and region functions
- [Examples](docs/examples.md) - Usage examples and complete code samples
