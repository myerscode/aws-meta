# Services Package

The `services` package provides functions for working with AWS services and regions, including listing services, finding services in specific regions, and getting service operations.

## Package Import

```go
import "github.com/myerscode/aws-meta/pkg/services"
```

## Functions

### ServiceMeta()

Returns detailed information about all AWS services including service IDs, full names, operations, and metadata.

```go
func ServiceMeta() (aws.ServiceSchemas, error)
```

**Returns:**
- `aws.ServiceSchemas`: A map of service schemas containing detailed service information (keyed by service ID)
- `error`: An error if there's an issue reading the manifest

**Example:**
```go
services, err := services.ServiceMeta()
if err != nil {
    log.Fatal(err)
}
for serviceId, service := range services {
    fmt.Printf("Service: %s\n", service.ServiceFullName)
    fmt.Printf("  ID: %s\n", service.ServiceId)
    fmt.Printf("  Operations: %d\n", len(service.Operations))
    fmt.Printf("  Protocol: %s\n", service.Protocol)
}
```

**Service Schema Fields:**
- `ServiceId`: Service identifier (e.g., "S3", "EC2", "Lambda")
- `ServiceFullName`: Full service name (e.g., "Amazon Simple Storage Service")
- `APIVersion`: API version for the service
- `EndpointPrefix`: Endpoint prefix used in service URLs
- `Protocol`: Communication protocol (e.g., "rest-xml", "json")
- `Operations`: Array of available operations for the service
- `GlobalEndpoint`: Global endpoint URL (if applicable)
- `SignatureVersion`: AWS signature version used
- `JSONVersion`: JSON version (for JSON-based services)
- `TargetPrefix`: Target prefix (for JSON-based services)

### ServiceMetaForRegion(regionName string) (aws.ServiceSchemas, error)

Returns detailed information about AWS services available in the specified region. This function combines the regional availability data with comprehensive service metadata.

```go
func ServiceMetaForRegion(regionName string) (aws.ServiceSchemas, error)
```

**Parameters:**
- `regionName`: The AWS region name (e.g., "us-east-1", "eu-west-1")

**Returns:**
- `aws.ServiceSchemas`: A map of service schemas containing detailed service information for services available in the region
- `error`: An error if the region name is invalid or if there's an issue reading the manifest

**Example:**
```go
services, err := services.ServiceMetaForRegion("us-east-1")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d services available in us-east-1\n", len(services))
for serviceId, service := range services {
    fmt.Printf("Service: %s\n", service.ServiceFullName)
    fmt.Printf("  ID: %s\n", service.ServiceId)
    fmt.Printf("  Endpoint: %s\n", service.EndpointPrefix)
    fmt.Printf("  Operations: %d\n", len(service.Operations))
}
```

**Notes:**
- The returned services are filtered to only include those available in the specified region
- Some endpoint prefixes may map to multiple services (e.g., "rds" includes RDS, Neptune, and DocumentDB)
- The function returns the same service metadata structure as `ServiceMeta()` but filtered by regional availability

**Error Cases:**
- Invalid region name: `invalid region name: invalid-region-123`
- Empty region name: `invalid region name: `
- Case sensitivity: Region names must match exactly (e.g., "us-east-1", not "US-EAST-1")

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

### AllServiceNames()

Returns a list of all AWS service full names.

```go
func AllServiceNames() []string
```

**Returns:**
- `[]string`: A sorted list of AWS service full names (e.g., "AWS Certificate Manager", "Amazon Simple Storage Service")

**Example:**
```go
serviceNames := services.AllServiceNames()
fmt.Printf("Total services: %d\n", len(serviceNames))
for _, service := range serviceNames[:3] { // Show first 3
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
serviceList, err := services.ServicesInRegion("us-east-1")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Services in us-east-1: %d\n", len(serviceList))
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

Service and region data is sourced from the official AWS Botocore library:
- **Regions**: From `botocore/data/endpoints.json`
- **Services**: From `botocore/data/*/service-*.json` files
- **Operations**: Extracted from individual service definition files

## Error Handling

Functions that return errors follow Go conventions:
- **Success**: `error` is `nil`
- **Failure**: `error` contains a descriptive message
- **Input validation**: Invalid inputs return specific error messages
- **Case sensitivity**: All inputs are case-sensitive and must match AWS conventions

## Performance Notes

- All data is loaded from pre-generated manifest files
- Functions return quickly as data is already parsed and sorted
- No network calls are made during function execution
- Data is cached in memory after first access