# Methods Reference

This document provides detailed information about all available methods in the AWS Metadata package.

## Partitions

### AllPartitionNames()

Returns a list of all AWS partition names.

```go
func AllPartitionNames() []string
```

**Returns:**
- `[]string`: A list of AWS partition names (e.g., "aws", "aws-cn", "aws-us-gov")

**Example:**
```go
partitions := services.AllPartitionNames()
fmt.Printf("Available partitions: %v\n", partitions)
// Output: Available partitions: [aws aws-cn aws-us-gov]
```

### ListPartitions()

Returns detailed information about all AWS partitions including DNS suffixes, region patterns, and region lists.

```go
func ListPartitions() (aws.PartitionSchemas, error)
```

**Package:** `github.com/myerscode/aws-meta/pkg/partitions`

**Returns:**
- `aws.PartitionSchemas`: A slice of partition schemas containing detailed partition information
- `error`: An error if there's an issue reading the manifest

**Example:**
```go
import "github.com/myerscode/aws-meta/pkg/partitions"

partitionList, err := partitions.ListPartitions()
if err != nil {
    log.Fatal(err)
}
for _, partition := range partitionList {
    fmt.Printf("Partition: %s\n", partition.ID)
    fmt.Printf("  DNS Suffix: %s\n", partition.DNSSuffix)
    fmt.Printf("  Regions: %d\n", len(partition.Regions))
}
```

**Partition Schema Fields:**
- `ID`: Partition identifier (e.g., "aws", "aws-cn")
- `DNSSuffix`: DNS suffix for the partition (e.g., "amazonaws.com")
- `DualStackDNSSuffix`: Dual-stack DNS suffix
- `RegionRegex`: Regular expression pattern for region names
- `ImplicitGlobalRegion`: Default global region for the partition
- `Regions`: Array of regions within the partition

**Partitions Include:**
- **aws**: Standard AWS partition (most regions)
- **aws-cn**: China partition (Beijing and Ningxia regions)
- **aws-us-gov**: AWS GovCloud partition (US government regions)
- **aws-eusc**: EU Sovereign Cloud partition
- **aws-iso**: US ISO partition (air-gapped)
- **aws-iso-b**: US ISOB partition (air-gapped)
- **aws-iso-e**: EU ISOE partition (air-gapped)
- **aws-iso-f**: US ISOF partition (air-gapped)

### CommercialPartitions()

Returns a list of AWS commercial partition names. Commercial partitions are the standard AWS partitions available to the general public.

```go
func CommercialPartitions() ([]string, error)
```

**Returns:**
- `[]string`: A list of commercial partition names (typically just "aws")
- `error`: An error if there's an issue reading the manifest

**Example:**
```go
partitions, err := services.CommercialPartitions()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Commercial partitions: %v\n", partitions)
// Output: Commercial partitions: [aws]
```

### SovereignPartitions()

Returns a list of AWS sovereign partition names. Sovereign partitions are operated by local entities to meet specific regulatory requirements.

```go
func SovereignPartitions() ([]string, error)
```

**Returns:**
- `[]string`: A list of sovereign partition names (e.g., "aws-cn", "aws-us-gov", "aws-eusc")
- `error`: An error if there's an issue reading the manifest

**Example:**
```go
partitions, err := services.SovereignPartitions()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Sovereign partitions: %v\n", partitions)
// Output: Sovereign partitions: [aws-cn aws-eusc aws-us-gov]
```

**Sovereign Partitions Include:**
- **aws-cn**: China partition operated by local Chinese entities
- **aws-us-gov**: US Government partition for federal agencies
- **aws-eusc**: EU Sovereign Cloud for European regulatory compliance

### IsolatedPartitions()

Returns a list of AWS isolated partition names. Isolated partitions are air-gapped environments for highly sensitive workloads.

```go
func IsolatedPartitions() ([]string, error)
```

**Returns:**
- `[]string`: A list of isolated partition names (e.g., "aws-iso", "aws-iso-b", "aws-iso-e", "aws-iso-f")
- `error`: An error if there's an issue reading the manifest

**Example:**
```go
partitions, err := services.IsolatedPartitions()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Isolated partitions: %v\n", partitions)
// Output: Isolated partitions: [aws-iso aws-iso-b aws-iso-e aws-iso-f]
```

**Isolated Partitions Include:**
- **aws-iso**: US ISO partition for intelligence community
- **aws-iso-b**: US ISOB partition for classified workloads
- **aws-iso-e**: EU ISOE partition for European classified workloads
- **aws-iso-f**: US ISOF partition for additional classified requirements

## Regions

### AllRegionNames()

Returns a list of all AWS region names across all partitions.

```go
func AllRegionNames() []string
```

**Returns:**
- `[]string`: A sorted list of all AWS region names (e.g., "us-east-1", "eu-west-1", "ap-southeast-1")

**Example:**
```go
regions := services.AllRegionNames()
fmt.Printf("Total regions: %d\n", len(regions))
for _, region := range regions[:5] { // Show first 5
    fmt.Printf("  - %s\n", region)
}
```

**Region Types Include:**
- **Standard regions**: us-east-1, us-west-2, eu-west-1, etc.
- **China regions**: cn-north-1, cn-northwest-1
- **GovCloud regions**: us-gov-east-1, us-gov-west-1
- **Opt-in regions**: Various newer regions that require explicit opt-in

## Services

### AllServiceNames()

Returns a list of all AWS service full names.

```go
func AllServiceNames() []string
```

**Returns:**
- `[]string`: A sorted list of AWS service full names (e.g., "AWS Certificate Manager", "Amazon Simple Storage Service")

**Example:**
```go
services := services.AllServiceNames()
fmt.Printf("Total services: %d\n", len(services))
for _, service := range services[:3] { // Show first 3
    fmt.Printf("  - %s\n", service)
}
// Output might include:
//   - AWS Certificate Manager
//   - Amazon Simple Storage Service
//   - Amazon Elastic Compute Cloud
```

### ServicesInRegion(regionName string) ([]string, error)

Returns a list of service names available in the specified region.

```go
func ServicesInRegion(regionName string) ([]string, error)
```

**Parameters:**
- `regionName`: The AWS region name (e.g., "us-east-1", "eu-west-1")

**Returns:**
- `[]string`: A sorted list of service names available in the region
- `error`: An error if the region name is invalid or if there's an issue reading the manifest

**Example:**
```go
services, err := services.ServicesInRegion("us-east-1")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Services in us-east-1: %d\n", len(services))
```

**Error Cases:**
- Invalid region name: `invalid region name: invalid-region-123`
- Empty region name: `invalid region name: `
- Case sensitivity: Region names must match exactly (e.g., "us-east-1", not "US-EAST-1")

### ServiceOperationsByServiceId(serviceId string) ([]string, error)

Returns a list of operation names available for the specified service ID.

```go
func ServiceOperationsByServiceId(serviceId string) ([]string, error)
```

**Parameters:**
- `serviceId`: The AWS service ID (e.g., "ACM", "S3", "EC2")

**Returns:**
- `[]string`: A sorted list of operation names available for the service
- `error`: An error if the service ID is invalid or if there's an issue reading the manifest

**Example:**
```go
operations, err := services.ServiceOperationsByServiceId("ACM")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("ACM operations: %d\n", len(operations))
// Might show: ACM operations: 16
```

**Common Service IDs:**
- `ACM`: AWS Certificate Manager
- `S3`: Amazon Simple Storage Service
- `EC2`: Amazon Elastic Compute Cloud
- `Lambda`: AWS Lambda
- `RDS`: Amazon Relational Database Service

**Error Cases:**
- Invalid service ID: `invalid service ID: INVALID-SERVICE`
- Case sensitivity: Service IDs are case-sensitive and typically uppercase

### ServiceOperationsByServiceName(serviceName string) ([]string, error)

Returns a list of operation names available for the specified service name.

```go
func ServiceOperationsByServiceName(serviceName string) ([]string, error)
```

**Parameters:**
- `serviceName`: The AWS service full name (e.g., "AWS Certificate Manager", "Amazon Simple Storage Service")

**Returns:**
- `[]string`: A sorted list of operation names available for the service
- `error`: An error if the service name is invalid or if there's an issue reading the manifest

**Example:**
```go
operations, err := services.ServiceOperationsByServiceName("AWS Certificate Manager")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("AWS Certificate Manager operations: %d\n", len(operations))
```

**Common Service Names:**
- `AWS Certificate Manager`: SSL/TLS certificate management
- `Amazon Simple Storage Service`: Object storage service
- `Amazon Elastic Compute Cloud`: Virtual server hosting
- `AWS Lambda`: Serverless compute service
- `Amazon Relational Database Service`: Managed database service

**Error Cases:**
- Invalid service name: `invalid service name: Invalid Service Name`
- Case sensitivity: Service names must match exactly as they appear in AWS documentation

### ServiceOperations(serviceId string) ([]string, error)

Alias for `ServiceOperationsByServiceId` for convenience.

```go
func ServiceOperations(serviceId string) ([]string, error)
```

This function provides the same functionality as `ServiceOperationsByServiceId` but with a shorter name for convenience.

**Parameters:**
- `serviceId`: The AWS service ID (e.g., "ACM", "S3", "EC2")

**Returns:**
- `[]string`: A sorted list of operation names available for the service
- `error`: An error if the service ID is invalid or if there's an issue reading the manifest

**Example:**
```go
operations, err := services.ServiceOperations("S3")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("S3 operations: %d\n", len(operations))
// Might show: S3 operations: 108
```

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