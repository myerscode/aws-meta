package aws

// Top-level structure for partitions.json
type EndpointFile struct {
	EndpointPartitions []EndpointFileEndpointPartitions `json:"partitions"`
}

// A single partition (e.g. "aws", "aws-cn")
type EndpointFileEndpointPartitions struct {
	ID           string                         `json:"partition"`
	RegionRegex  string                         `json:"regionRegex"`
	Default      EndpointFileDefaults           `json:"defaults"`
	DNSSuffix    string                         `json:"dnsSuffix"`
	SupportsFIPS bool                           `json:"supportsFIPS"`
	SupportsDual bool                           `json:"supportsDualStack"`
	Regions      map[string]EndpointFileRegion  `json:"regions"`
	Services     map[string]EndpointFileService `json:"services"`
	Outputs      EndpointFileOutputs            `json:"outputs"`
}

// Fields under "outputs"
type EndpointFileOutputs struct {
	DNSSuffix            string `json:"dnsSuffix"`
	DualStackDNSSuffix   string `json:"dualStackDnsSuffix,omitempty"`
	ImplicitGlobalRegion string `json_json:"implicitGlobalRegion,omitempty"`
	Name                 string `json:"name"`
	SupportsDualStack    bool   `json:"supportsDualStack,omitempty"`
	SupportsFIPS         bool   `json:"supportsFIPS,omitempty"`
}

type EndpointFileRegion struct {
	Description string `json:"description,omitempty"`
	OptInStatus string `json:"optInStatus,omitempty"`
}

type EndpointFileService struct {
	IsRegionalized    bool                            `json:"isRegionalized,omitempty"`
	PartitionEndpoint string                          `json:"partitionEndpoint,omitempty"`
	Defaults          EndpointFileDefaults            `json:"defaults,omitempty"`
	Endpoints         map[string]EndpointFileDefaults `json:"endpoints"`
}

type EndpointFileDefaults struct {
	Hostname          string   `json:"hostname,omitempty"`
	Protocols         []string `json:"protocols,omitempty"`
	SignatureVersions []string `json:"signatureVersions,omitempty"`
	CredentialScope   struct {
		Region  string `json:"region,omitempty"`
		Service string `json:"service,omitempty"`
	} `json:"credentialScope,omitempty"`
	SSLCommonName string `json:"sslCommonName,omitempty"`
	Variants      []struct {
		Hostname string   `json:"hostname"`
		Tags     []string `json:"tags"`
	} `json:"variants,omitempty"`
}
