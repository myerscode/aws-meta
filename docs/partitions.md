# Partitions Package

The `partitions` package provides functions for working with AWS partitions, including listing all partitions and categorizing them by type.

## Package Import

```go
import "github.com/myerscode/aws-meta/pkg/partitions"
```

## Functions

### AllPartitionNames()

Returns a list of all AWS partition names.

```go
func AllPartitionNames() []string
```

**Returns:**
- `[]string`: A list of AWS partition names (e.g., "aws", "aws-cn", "aws-us-gov")

**Example:**
```go
partitionNames := partitions.AllPartitionNames()
fmt.Printf("Available partitions: %v\n", partitionNames)
// Output: Available partitions: [aws aws-cn aws-us-gov]
```

### List()

Returns detailed information about all AWS partitions including DNS suffixes, region patterns, and region lists.

```go
func List() (aws.PartitionSchemas, error)
```

**Returns:**
- `aws.PartitionSchemas`: A slice of partition schemas containing detailed partition information
- `error`: An error if there's an issue reading the manifest

**Example:**
```go
partitionList, err := partitions.List()
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
partitionList, err := partitions.CommercialPartitions()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Commercial partitions: %v\n", partitionList)
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
partitionList, err := partitions.SovereignPartitions()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Sovereign partitions: %v\n", partitionList)
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
partitionList, err := partitions.IsolatedPartitions()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Isolated partitions: %v\n", partitionList)
// Output: Isolated partitions: [aws-iso aws-iso-b aws-iso-e aws-iso-f]
```

**Isolated Partitions Include:**
- **aws-iso**: US ISO partition for intelligence community
- **aws-iso-b**: US ISOB partition for classified workloads
- **aws-iso-e**: EU ISOE partition for European classified workloads
- **aws-iso-f**: US ISOF partition for additional classified requirements

## Data Sources

Partition data is sourced from the official AWS Botocore library:
- **Partitions**: From `botocore/data/partitions.json`

## Error Handling

Functions that return errors follow Go conventions:
- **Success**: `error` is `nil`
- **Failure**: `error` contains a descriptive message
- **Input validation**: Invalid inputs return specific error messages

## Performance Notes

- All data is loaded from pre-generated manifest files
- Functions return quickly as data is already parsed and sorted
- No network calls are made during function execution
- Data is cached in memory after first access