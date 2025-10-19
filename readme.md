# AWS Metadata
> A Go package for exposing information about AWS Partitions, Regions and Services

[![Go](https://github.com/myerscode/aws-meta/actions/workflows/go.yml/badge.svg)](https://github.com/myerscode/aws-meta/actions/workflows/go.yml)
[![Nightly Release](https://github.com/myerscode/aws-meta/actions/workflows/release.yml/badge.svg)](https://github.com/myerscode/aws-meta/actions/workflows/release.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/myerscode/aws-meta.svg)](https://pkg.go.dev/github.com/myerscode/aws-meta)
[![Go Report Card](https://goreportcard.com/badge/github.com/myerscode/aws-meta)](https://goreportcard.com/report/github.com/myerscode/aws-meta)

Knowing what AWS services are available in which regions and partitions can be a bit of a pain; 
this package aims to keep an up to date reference for the available endpoints for all AWS services.

## Features

- **Partitions**: Get all AWS partition names (aws, aws-cn, aws-us-gov)
- **Regions**: List all AWS regions across all partitions
- **Services**: Access AWS service information and operations
- **Regional Services**: Find which services are available in specific regions
- **Service Operations**: Get all available operations for any AWS service

## Installation

```bash
go get github.com/myerscode/aws-meta
```

## Quick Start

### Working with Partitions
```go
import "github.com/myerscode/aws-meta/pkg/partitions"

// Get all partition names
partitions := partitions.AllPartitionNames()
fmt.Printf("Available partitions: %v\n", partitions)

// Get detailed partition information
partitionList, err := partitions.List()
if err != nil {
    log.Fatal(err)
}

// Get partitions by type
commercial, _ := partitions.CommercialPartitions()
sovereign, _ := partitions.SovereignPartitions()
isolated, _ := partitions.IsolatedPartitions()
```

### Working with Services and Regions
```go
import "github.com/myerscode/aws-meta/pkg/services"

// Get all regions
regions := services.AllRegionNames()

// Get detailed service metadata
services, err := services.ServiceMeta()
if err != nil {
    log.Fatal(err)
}

// Get services available in a specific region
regionServices, err := services.ServiceMetaForRegion("us-east-1")
if err != nil {
    log.Fatal(err)
}

// Get service operations
operations, err := services.ServiceOperations("S3")
if err != nil {
    log.Fatal(err)
}
```

## Documentation

- [API Reference](docs/methods.md) - Overview of all packages and functions
- [Partitions Package](docs/partitions.md) - Partition-related functions
- [Services Package](docs/services.md) - Service and region functions
- [Examples](docs/examples.md) - Usage examples and complete code samples
