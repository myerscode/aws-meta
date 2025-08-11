# AWS Metadata
> A package for exposing information about AWS Partitions, Regions and Services


Knowing what AWS services are available in which regions and partitions can be a bit of a pain; 
this app aims to keep an up to date reference for the available endpoints for all AWS services.

## Partitions

```go
// Returns a list of all AWS partition names
AllPartitionNames()
```

### Isolated

### Sovereign

## Regions

```go
// Returns a list of all AWS region names
AllRegionNames()
```

### Regions For Partition

## Services

```go
// Returns a list of all AWS service names
AllServiceNames()

// Returns a list of service names available in the specified region
// Returns an error if the region name is invalid
ServicesInRegion(regionName string) ([]string, error)
```

### Services In Region

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

## Examples

See the `examples/` directory for working code examples:

- `examples/list_all.go` - General usage examples
- `examples/services_in_region.go` - Demonstrates the `ServicesInRegion` function

Run any example with:
```bash
go run examples/example_name.go
```
