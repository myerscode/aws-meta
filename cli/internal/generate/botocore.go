package generate

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/myerscode/aws-meta/cli/internal/github"
	"github.com/myerscode/aws-meta/cli/internal/util"
	"github.com/myerscode/aws-meta/internal/aws"
)

type Botocore struct {
	Repo github.Repo
	// OutputDir is the directory the generated pkg/data files are written
	// beneath, typically the repository root.
	OutputDir string
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

// BotoPartitionsFiles represents the partition.json file in botocore
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

// EndpointFile is the top-level structure for endpoints.json
type EndpointFile struct {
	EndpointPartitions []EndpointFilePartition `json:"partitions"`
}

type EndpointFilePartition struct {
	ID       string                         `json:"partition"`
	Regions  map[string]EndpointFileRegion  `json:"regions"`
	Services map[string]EndpointFileService `json:"services"`
}

type EndpointFileRegion struct {
	Description string `json:"description,omitempty"`
}

type EndpointFileService struct {
	Endpoints map[string]json.RawMessage `json:"endpoints"`
}

// BotoServiceDataSources maps service names to their data source info.
type BotoServiceDataSources map[string]BotoDataSource

type BotoDataSource struct {
	ApiVersion string
	Filename   string
}

// DownloadTagData fetches all needed botocore data for a tag in a single
// tarball download, replacing hundreds of individual API calls with one
// HTTP request.
func (bc Botocore) DownloadTagData(tag github.RepoTag) (github.TarballFiles, error) {
	pathPrefixes := []string{
		"botocore/data/",
	}
	return bc.Repo.DownloadAndExtract(tag, pathPrefixes)
}

// GeneratePartitionList builds partition metadata from pre-extracted tarball data.
func (bc Botocore) GeneratePartitionList(tag github.RepoTag, files github.TarballFiles) aws.PartitionSchemas {
	meta, metaErr := parsePartitionMeta(files)
	if metaErr != nil {
		util.PrintErrorAndExit(metaErr)
	}

	fmt.Printf("Version: %s\n", meta.Version)
	fmt.Printf("Number of partitions: %d\n", len(meta.Partitions))

	var partitionSchemas aws.PartitionSchemas

	for _, partition := range meta.Partitions {
		var partitionRegions []aws.PartitionRegion

		for regionID, region := range partition.Regions {
			partitionRegions = append(partitionRegions, aws.PartitionRegion{
				RegionId:   regionID,
				RegionName: region.Description,
			})
		}

		err := util.SortByField(&partitionRegions, "RegionId")
		if err != nil {
			util.PrintErrorAndExit(err)
		}

		partitionSchemas = append(partitionSchemas, aws.PartitionSchema{
			ID:                   partition.ID,
			RegionRegex:          partition.RegionRegex,
			Regions:              partitionRegions,
			DNSSuffix:            partition.Outputs.DNSSuffix,
			DualStackDNSSuffix:   partition.Outputs.DualStackDNSSuffix,
			ImplicitGlobalRegion: partition.Outputs.ImplicitGlobalRegion,
		})
	}

	err := util.SortByField(&partitionSchemas, "ID")
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	err = SaveManifestFile(bc.OutputDir, partitionSchemas, "botocore.partitions.json")
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	err = SaveArchiveFile(bc.OutputDir, partitionSchemas, fmt.Sprintf("botocore.partitions.%s.json", tag.Name))
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	return partitionSchemas
}

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

// serviceSchemaRe matches service schema file paths in the tarball.
var serviceSchemaRe = regexp.MustCompile(`^botocore/data/(?P<service>.+?)/(?P<apiVersion>.+?)/service-\d+\.json$`)

// GenerateServiceList builds service metadata from pre-extracted tarball data.
func (bc Botocore) GenerateServiceList(tag github.RepoTag, files github.TarballFiles) aws.ServiceSchemas {
	util.LogInfo(fmt.Sprintf("Generating service list from tarball data for tag: %s", tag.Name))

	dataSources := findServiceDataSources(files)
	serviceSchemas := aws.ServiceSchemas{}

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

		serviceSchema := aws.ServiceSchema{
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

	err := SaveManifestFile(bc.OutputDir, serviceSchemas, "botocore.services.json")
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	err = SaveArchiveFile(bc.OutputDir, serviceSchemas, fmt.Sprintf("botocore.services.%s.json", tag.Name))
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	return serviceSchemas
}

func findServiceDataSources(files github.TarballFiles) BotoServiceDataSources {
	dataSourceMap := make(BotoServiceDataSources)

	for path := range files {
		matches := serviceSchemaRe.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		service := matches[serviceSchemaRe.SubexpIndex("service")]
		apiVersion := matches[serviceSchemaRe.SubexpIndex("apiVersion")]

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
func (bc Botocore) GenerateRegionServicesList(tag github.RepoTag, files github.TarballFiles) aws.RegionSchemas {
	util.LogInfo(fmt.Sprintf("GenerateRegionServicesList from tarball data for tag: %s", tag.Name))

	endpointData, endpointDataError := parseEndpointData(files)
	if endpointDataError != nil {
		util.PrintErrorAndExit(endpointDataError)
	}

	var summaries aws.RegionSchemas

	for _, partition := range endpointData.EndpointPartitions {
		summary := aws.RegionSchema{
			PartitionID: partition.ID,
			Regions:     []aws.RegionSummary{},
		}

		for regionName := range partition.Regions {
			var servicesInRegion []string
			for serviceName, service := range partition.Services {
				if _, ok := service.Endpoints[regionName]; ok {
					servicesInRegion = append(servicesInRegion, serviceName)
				}
			}
			summary.Regions = append(summary.Regions, aws.RegionSummary{
				RegionName: regionName,
				Services:   servicesInRegion,
			})
		}

		summaries = append(summaries, summary)
	}

	sortRegionSchemas(summaries)

	err := SaveManifestFile(bc.OutputDir, summaries, "botocore.regions.json")
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	err = SaveArchiveFile(bc.OutputDir, summaries, fmt.Sprintf("botocore.regions.%s.json", tag.Name))
	if err != nil {
		util.PrintErrorAndExit(err)
	}

	return summaries
}

func parseEndpointData(files github.TarballFiles) (EndpointFile, error) {
	var endpointFile EndpointFile

	var blob []byte
	var ok bool

	blob, ok = files["botocore/data/endpoints.json"]
	if !ok {
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

func sortRegionSchemas(schemas aws.RegionSchemas) {
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].PartitionID < schemas[j].PartitionID
	})

	for i := range schemas {
		sort.Slice(schemas[i].Regions, func(a, b int) bool {
			return schemas[i].Regions[a].RegionName < schemas[i].Regions[b].RegionName
		})

		for j := range schemas[i].Regions {
			sort.Strings(schemas[i].Regions[j].Services)
		}
	}
}
