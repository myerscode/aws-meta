package aws

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/myerscode/aws-meta/internal/github"
	"github.com/myerscode/aws-meta/internal/util"
)

type Botocore struct {
	Repo github.Repo
}

// DataSchema represents a service schema file
type DataSchema struct {
	Version    string                `json:"version"`
	Metadata   MetaData              `json:"metadata"`
	Operations map[string]Operations `json:"operations"`
}
type MetaData struct {
	ServiceId        string `json:"serviceId"`
	ServiceFullName  string `json:"serviceFullName"`
	EndpointPrefix   string `json:"endpointPrefix"`
	GlobalEndpoint   string `json:"globalEndpoint,omitempty"`
	SignatureVersion string `json:"signatureVersion,omitempty"`
	Protocol         string `json:"protocol,omitempty"`
	ApiVersion       string `json:"apiVersion"`
	JSONVersion      string `json:"jsonVersion,omitempty"`
	TargetPrefix     string `json:"targetPrefix,omitempty"`
}
type Operations struct {
	Name string `json:"name"`
}

type PartitionSchemas []PartitionSchema

type PartitionSchema struct {
	ID                   string
	RegionRegex          string
	DNSSuffix            string
	DualStackDNSSuffix   string
	ImplicitGlobalRegion string
	Regions              []PartitionRegion
}

type PartitionRegion struct {
	RegionId   string
	RegionName string
}

// DownloadTagData fetches all needed botocore data for a tag in a single
// tarball download, replacing hundreds of individual API calls with one
// HTTP request.
func (bc Botocore) DownloadTagData(tag github.RepoTag) (github.TarballFiles, error) {
	// Only extract the paths we actually need from the ~50MB tarball
	pathPrefixes := []string{
		"botocore/data/partitions.json",
		"botocore/data/endpoints.json",
		"botocore/data/", // service schema files live under botocore/data/{service}/{version}/
	}

	return bc.Repo.DownloadAndExtract(tag, pathPrefixes)
}

// GeneratePartitionList builds partition metadata from pre-extracted tarball data.
func (bc Botocore) GeneratePartitionList(tag github.RepoTag, files github.TarballFiles) PartitionSchemas {
	meta, metaErr := parsePartitionMeta(files)
	if metaErr != nil {
		util.PrintErrorAndExit(metaErr)
	}

	fmt.Printf("Version: %s\n", meta.Version)
	fmt.Printf("Number of partitions: %d\n", len(meta.Partitions))

	var partitionSchemas PartitionSchemas

	if len(meta.Partitions) > 0 {
		for _, partition := range meta.Partitions {
			var partitionRegions []PartitionRegion

			for regionID, region := range partition.Regions {
				partitionRegions = append(partitionRegions, PartitionRegion{
					RegionId:   regionID,
					RegionName: region.Description,
				})
			}

			err := util.SortByField(&partitionRegions, "RegionId")
			if err != nil {
				util.PrintErrorAndExit(err)
			}

			partitionSchemas = append(partitionSchemas, PartitionSchema{
				ID:                   partition.ID,
				RegionRegex:          partition.RegionRegex,
				Regions:              partitionRegions,
				DNSSuffix:            partition.Outputs.DNSSuffix,
				DualStackDNSSuffix:   partition.Outputs.DualStackDNSSuffix,
				ImplicitGlobalRegion: partition.Outputs.ImplicitGlobalRegion,
			})
		}
	}

	err := util.SortByField(&partitionSchemas, "ID")
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	err = SaveManifestFile(partitionSchemas, "botocore.partitions.json")
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	err = SaveArchiveFile(partitionSchemas, fmt.Sprintf("botocore.partitions.%s.json", tag.Name))
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	return partitionSchemas
}

// parsePartitionMeta reads partition data from the pre-extracted files map.
func parsePartitionMeta(files github.TarballFiles) (BotoPartitionsFiles, error) {
	var partition BotoPartitionsFiles

	blob, ok := files["botocore/data/partitions.json"]
	if !ok {
		return BotoPartitionsFiles{}, fmt.Errorf("partitions.json not found in tarball")
	}

	if err := json.Unmarshal(blob, &partition); err != nil {
		return BotoPartitionsFiles{}, err
	}

	return partition, nil
}

// GenerateServiceList builds service metadata from pre-extracted tarball data.
// Instead of making 300+ individual HTTP requests, it reads service schema
// files directly from the extracted tarball contents.
func (bc Botocore) GenerateServiceList(tag github.RepoTag, files github.TarballFiles) ServiceSchemas {
	util.LogInfo(fmt.Sprintf("Generating service list from tarball data for tag: %s", tag.Name))

	dataSources := findServiceDataSources(files)

	serviceSchemas := ServiceSchemas{}

	for _, dataSource := range dataSources {
		blob, ok := files[dataSource.Filename]
		if !ok {
			util.LogWarning(fmt.Sprintf("Service file not found in tarball: %s", dataSource.Filename))
			continue
		}

		dataSchema := DataSchema{}
		if err := json.Unmarshal(blob, &dataSchema); err != nil {
			util.LogWarning(fmt.Sprintf("Failed to parse service schema %s: %v", dataSource.Filename, err))
			continue
		}

		var operations []string
		for operation := range dataSchema.Operations {
			operations = append(operations, operation)
		}

		serviceSchema := ServiceSchema{
			APIVersion:       dataSource.ApiVersion,
			ServiceId:        dataSchema.Metadata.ServiceId,
			ServiceFullName:  dataSchema.Metadata.ServiceFullName,
			EndpointPrefix:   dataSchema.Metadata.EndpointPrefix,
			GlobalEndpoint:   dataSchema.Metadata.GlobalEndpoint,
			SignatureVersion: dataSchema.Metadata.SignatureVersion,
			Protocol:         dataSchema.Metadata.Protocol,
			JSONVersion:      dataSchema.Metadata.JSONVersion,
			TargetPrefix:     dataSchema.Metadata.TargetPrefix,
			Operations:       util.Sort(operations),
		}

		serviceSchemas[serviceSchema.ServiceId] = serviceSchema
	}

	err := SaveManifestFile(serviceSchemas, "botocore.services.json")
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	err = SaveArchiveFile(serviceSchemas, fmt.Sprintf("botocore.services.%s.json", tag.Name))
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	return serviceSchemas
}

// findServiceDataSources discovers service schema files from the extracted
// tarball file map, picking the latest API version for each service.
// This replaces the tree API call + regex matching over the API response.
func findServiceDataSources(files github.TarballFiles) BotoServiceDataSources {
	dataSourceMap := make(BotoServiceDataSources)

	re := regexp.MustCompile(`^botocore/data/(?P<service>.+?)/(?P<apiVersion>.+?)/service-\d+\.json$`)

	for path := range files {
		matches := re.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		service := matches[re.SubexpIndex("service")]
		apiVersion := matches[re.SubexpIndex("apiVersion")]

		if existing, ok := dataSourceMap[service]; ok {
			if apiVersion < existing.ApiVersion {
				continue
			}
		}

		dataSourceMap[service] = BotoDataSource{
			ApiVersion: apiVersion,
			Filename:   path,
		}
	}

	return dataSourceMap
}

// GenerateRegionServicesList builds region-to-service mappings from pre-extracted tarball data.
func (bc Botocore) GenerateRegionServicesList(tag github.RepoTag, files github.TarballFiles) RegionSchemas {
	util.LogInfo(fmt.Sprintf("GenerateRegionServicesList from tarball data for tag: %s", tag.Name))

	endpointData, endpointDataError := parseEndpointData(files)
	if endpointDataError != nil {
		util.PrintErrorAndExit(endpointDataError)
	}

	var summaries RegionSchemas

	for _, partition := range endpointData.EndpointPartitions {
		summary := RegionSchema{
			PartitionID: partition.ID,
			Regions:     []RegionSummary{},
		}

		for regionName := range partition.Regions {
			var servicesInRegion []string
			for serviceName, service := range partition.Services {
				if _, ok := service.Endpoints[regionName]; ok {
					servicesInRegion = append(servicesInRegion, serviceName)
				}
			}
			summary.Regions = append(summary.Regions, RegionSummary{
				RegionName: regionName,
				Services:   servicesInRegion,
			})
		}

		summaries = append(summaries, summary)
	}

	sortRegionSchemas(summaries)

	err := SaveManifestFile(summaries, "botocore.regions.json")
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	err = SaveArchiveFile(summaries, fmt.Sprintf("botocore.regions.%s.json", tag.Name))
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	return summaries
}

// parseEndpointData reads endpoint data from the pre-extracted files map.
func parseEndpointData(files github.TarballFiles) (EndpointFile, error) {
	var endpointFile EndpointFile

	// Try both possible locations for the endpoints file
	var blob []byte
	var ok bool

	blob, ok = files["botocore/data/endpoints.json"]
	if !ok {
		// Some versions might not have the nested path
		for path, data := range files {
			if strings.HasSuffix(path, "endpoints.json") {
				blob = data
				ok = true
				break
			}
		}
	}

	if !ok {
		return EndpointFile{}, fmt.Errorf("endpoints.json not found in tarball")
	}

	if err := json.Unmarshal(blob, &endpointFile); err != nil {
		return EndpointFile{}, err
	}

	return endpointFile, nil
}

func sortRegionSchemas(schemas RegionSchemas) {
	// Sort by PartitionID
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].PartitionID < schemas[j].PartitionID
	})

	for i := range schemas {
		// Sort each RegionSummary by RegionName
		sort.Slice(schemas[i].Regions, func(a, b int) bool {
			return schemas[i].Regions[a].RegionName < schemas[i].Regions[b].RegionName
		})

		for j := range schemas[i].Regions {
			// Sort services alphabetically
			sort.Strings(schemas[i].Regions[j].Services)
		}
	}
}
