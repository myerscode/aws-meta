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

Run any example with:
```bash
go run examples/example_name.go
```

## Complete Example

Here's a complete example showing how to use multiple functions together:

```go
package main

import (
    "fmt"
    "log"

    "github.com/myerscode/aws-meta/pkg/services"
)

func main() {
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