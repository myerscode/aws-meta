package aws

// BotoPartitionsFiles represents the partition.json file in boto
type BotoPartitionsFiles struct {
	Partitions []BotoPartitionsFilePartition `json:"partitions"`
	Version    string                        `json:"version"`
}

type BotoPartitionsFilePartition struct {
	ID      string `json:"id"`
	Outputs struct {
		DNSSuffix            string `json:"dnsSuffix"`
		DualStackDNSSuffix   string `json:"dualStackDnsSuffix"`
		ImplicitGlobalRegion string `json:"implicitGlobalRegion"`
		Name                 string `json:"name"`
		SupportsDualStack    bool   `json:"supportsDualStack"`
		SupportsFIPS         bool   `json:"supportsFIPS"`
	} `json:"outputs"`
	RegionRegex string                                        `json:"regionRegex"`
	Regions     map[string]BotoPartitionsFileRegionDefinition `json:"regions"`
}
type BotoPartitionsFileRegionDefinition struct {
	Description string `json:"description"`
}
