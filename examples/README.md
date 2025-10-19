# Examples

This directory contains example programs demonstrating how to use the aws-meta library.

## Running Examples

Each example is in its own subdirectory and can be run independently:

```bash
# List all services, regions, and partitions
go run examples/list-all/main.go

# Show detailed partition information
go run examples/list-partitions/main.go

# Show partition types (commercial, sovereign, isolated)
go run examples/partition-types/main.go

# Demonstrate service operations lookup
go run examples/service-operations/main.go

# Show services available in specific regions
go run examples/services-in-region/main.go
```

## Examples Overview

- **list-all**: Basic usage showing how to get all service names, region names, and partition names
- **list-partitions**: Detailed information about AWS partitions including regions and DNS suffixes
- **partition-types**: Demonstrates categorization of partitions into commercial, sovereign, and isolated types
- **service-operations**: Shows how to look up operations for specific AWS services
- **services-in-region**: Demonstrates finding which services are available in specific regions

## Shared Utilities

The `shared` package contains common utility functions used across examples to avoid code duplication.