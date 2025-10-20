# Regions Package

**Import:** `github.com/myerscode/aws-meta/pkg/regions`

The regions package provides comprehensive information about AWS regions across all partitions, including partition mapping and service availability.

## Functions

### ListAllRegions()

Returns detailed information about all AWS regions across all partitions including partition mapping and service availability.

```go
func ListAllRegions() ([]RegionInfo, error)
```

**Returns:**
- `[]RegionInfo`: A slice of region information containing detailed region data
- `error`: An error if there's an issue reading the manifest

**Example:**
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
}
```

## Data Structures

### RegionInfo

The `RegionInfo` struct contains comprehensive information about each AWS region:

```go
type RegionInfo struct {
    RegionId    string   `json:"regionId"`    // Region identifier (e.g., "us-east-1")
    RegionName  string   `json:"regionName"`  // Human-readable name (e.g., "US East (N. Virginia)")
    PartitionID string   `json:"partitionId"` // Partition the region belongs to (e.g., "aws")
    Services    []string `json:"services"`    // Array of service names available in the region
}
```

**Fields:**
- **RegionId**: Region identifier (e.g., "us-east-1", "eu-west-1")
- **RegionName**: Human-readable region name (e.g., "US East (N. Virginia)")
- **PartitionID**: Partition the region belongs to (e.g., "aws", "aws-cn")
- **Services**: Array of service names available in the region


## Usage Examples

### Basic Region Listing

```go
regions, err := regions.ListAllRegions()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d regions\n", len(regions))
for _, region := range regions {
    fmt.Printf("- %s (%s) in %s partition\n", 
        region.RegionId, region.RegionName, region.PartitionID)
}
```

### Group Regions by Partition

```go
regions, err := regions.ListAllRegions()
if err != nil {
    log.Fatal(err)
}

partitionRegions := make(map[string][]regions.RegionInfo)
for _, region := range regions {
    partitionRegions[region.PartitionID] = append(partitionRegions[region.PartitionID], region)
}

for partition, regions := range partitionRegions {
    fmt.Printf("Partition %s: %d regions\n", partition, len(regions))
}
```

### Find Regions with Most Services

```go
regions, err := regions.ListAllRegions()
if err != nil {
    log.Fatal(err)
}

var maxServices int
var topRegion regions.RegionInfo

for _, region := range regions {
    if len(region.Services) > maxServices {
        maxServices = len(region.Services)
        topRegion = region
    }
}

fmt.Printf("Region with most services: %s (%s) - %d services\n", 
    topRegion.RegionId, topRegion.RegionName, maxServices)
```

### Filter by Geographic Area

```go
regions, err := regions.ListAllRegions()
if err != nil {
    log.Fatal(err)
}

// Find all European regions in AWS partition
var europeanRegions []regions.RegionInfo
for _, region := range regions {
    if region.PartitionID == "aws" && strings.HasPrefix(region.RegionId, "eu-") {
        europeanRegions = append(europeanRegions, region)
    }
}

fmt.Printf("European regions: %d\n", len(europeanRegions))
for _, region := range europeanRegions {
    fmt.Printf("- %s (%s)\n", region.RegionId, region.RegionName)
}
```

## Error Handling

The `ListAllRegions()` function returns an error if:
- There's an issue reading the partition manifest
- There's an issue reading the region manifest
- The manifest files are corrupted or missing

```go
regions, err := regions.ListAllRegions()
if err != nil {
    log.Fatalf("Failed to load regions: %v", err)
}
```

## Performance Notes

- Data is loaded from pre-generated manifest files
- Function returns quickly as data is already parsed and sorted
- No network calls are made during function execution
- Data is cached in memory after first access
- All regions and their services are loaded in a single call

## Integration with Other Packages

The regions package complements other packages in the library:

```go
import (
    "github.com/myerscode/aws-meta/pkg/regions"
    "github.com/myerscode/aws-meta/pkg/services"
    "github.com/myerscode/aws-meta/pkg/partitions"
)

// Get detailed region information
regionList, _ := regions.ListAllRegions()

// Cross-reference with partition details
partitionList, _ := partitions.List()

// Verify service availability
servicesInRegion, _ := services.ServicesInRegion("us-east-1")
```

## Complete Example

See `examples/list_regions.go` for a complete, runnable example that demonstrates:
- Listing all regions with service counts
- Grouping regions by partition
- Finding regions with most/least services
- Geographic distribution analysis
- Regional statistics and comparisons