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
```
