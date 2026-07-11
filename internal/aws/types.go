package aws

// ServiceSchemas maps service IDs to their schema metadata.
type ServiceSchemas map[string]ServiceSchema

// ServiceSchema holds metadata for a single AWS service.
type ServiceSchema struct {
	APIVersion       string
	ServiceId        string
	ServiceFullName  string
	EndpointPrefix   string
	GlobalEndpoint   string   `json:"GlobalEndpoint,omitempty"`
	SignatureVersion string   `json:"SignatureVersion,omitempty"`
	Protocol         string   `json:"Protocol,omitempty"`
	JSONVersion      string   `json:"JSONVersion,omitempty"`
	TargetPrefix     string   `json:"TargetPrefix,omitempty"`
	Operations       []string `json:"Operations"`
}

// RegionSchemas holds region summaries grouped by partition.
type RegionSchemas []RegionSchema

// RegionSchema holds region summaries per partition.
type RegionSchema struct {
	PartitionID string          `json:"partition"`
	Regions     []RegionSummary `json:"regions"`
}

// RegionSummary holds services for a given region.
type RegionSummary struct {
	RegionName string   `json:"region"`
	Services   []string `json:"services"`
}

// PartitionSchemas is a slice of partition metadata.
type PartitionSchemas []PartitionSchema

// PartitionSchema holds metadata for a single AWS partition.
type PartitionSchema struct {
	ID                   string
	RegionRegex          string
	DNSSuffix            string
	DualStackDNSSuffix   string
	ImplicitGlobalRegion string
	Regions              []PartitionRegion
}

// PartitionRegion represents a region within a partition.
type PartitionRegion struct {
	RegionId   string
	RegionName string
}
