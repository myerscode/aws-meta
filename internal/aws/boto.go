package aws

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"sync"

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

func (bc Botocore) GeneratePartitionList(tag github.RepoTag) PartitionSchemas {

	meta, metaErr := bc.getPartitionMeta(tag)

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

func (bc Botocore) getPartitionMeta(tag github.RepoTag) (BotoPartitionsFiles, error) {
	var partition BotoPartitionsFiles

	blob, err := bc.Repo.GetBlobFromTag(tag.Name, "data/partitions.json")

	if err != nil {
		return BotoPartitionsFiles{}, err
	}

	if err := json.Unmarshal(blob, &partition); err != nil {
		return BotoPartitionsFiles{}, err
	}

	return partition, nil
}

func (bc Botocore) GenerateServiceList(tag github.RepoTag) ServiceSchemas {

	util.LogInfo(fmt.Sprintf("GetServiceDataSources for Tag: %s", tag.Name))

	dataSources, err := bc.getServiceDataSources(tag)

	if err != nil {
		util.PrintErrorAndExit(err)
	}

	var wg sync.WaitGroup
	serviceSchemaChannel := make(chan ServiceSchema, len(dataSources))

	for _, dataSource := range dataSources {
		wg.Add(1)
		go generateServiceSchema(&wg, bc, tag, dataSource, serviceSchemaChannel)
	}

	wg.Wait()

	close(serviceSchemaChannel)

	serviceSchemas := ServiceSchemas{}

	for serviceSchema := range serviceSchemaChannel {
		serviceSchemas[serviceSchema.ServiceId] = serviceSchema
	}

	err = SaveManifestFile(serviceSchemas, "botocore.services.json")

	if err != nil {
		util.PrintErrorAndExit(err)
	}

	err = SaveArchiveFile(serviceSchemas, fmt.Sprintf("botocore.services.%s.json", tag.Name))

	if err != nil {
		util.PrintErrorAndExit(err)
	}

	return serviceSchemas
}

func generateServiceSchema(wg *sync.WaitGroup, bc Botocore, tag github.RepoTag, dataSource BotoDataSource, serviceSchemaChan chan<- ServiceSchema) {
	defer wg.Done()

	dataSchema := DataSchema{}

	rawData, err := github.GetGithubRepoBlobs(bc.Repo.Owner, bc.Repo.RepoName, tag.Name, dataSource.Filename)

	if err != nil {
		util.PrintErrorAndExit(err)
	}

	if err := json.Unmarshal(rawData, &dataSchema); err != nil {
		util.PrintErrorAndExit(err)
	}

	var operations []string

	for operation := range dataSchema.Operations {
		operations = append(operations, operation)
	}

	// Create service schema with required fields
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

	serviceSchemaChan <- serviceSchema
}

func (bc Botocore) getServiceDataSources(tag github.RepoTag) (BotoServiceDataSources, error) {

	trees, err := bc.Repo.GetGithubRepoTrees(tag.Commit.SHA, "botocore/data")

	if err != nil {
		return map[string]BotoDataSource{}, err
	}

	dataSourceMap := map[string]BotoDataSource{}

	re := regexp.MustCompile(`(?P<service>.+?)/(?P<apiVersion>.+?)/service-\d.json`)

	for _, value := range trees {
		matches := re.FindStringSubmatch(value.Path)
		if matches == nil {
			continue
		}

		service := matches[re.SubexpIndex("service")]
		apiVersion := matches[re.SubexpIndex("apiVersion")]

		if _, ok := dataSourceMap[service]; ok {
			if apiVersion < dataSourceMap[service].ApiVersion {
				continue
			}
		}

		dataSourceMap[service] = BotoDataSource{
			ApiVersion: apiVersion,
			Filename:   fmt.Sprintf("%s/%s", "botocore/data", matches[0]),
			Sha:        value.Sha,
		}
	}

	return dataSourceMap, nil
}

func (bc Botocore) GenerateRegionServicesList(tag github.RepoTag) RegionSchemas {

	util.LogInfo(fmt.Sprintf("GenerateRegionServicesList for Tag: %s", tag.Name))

	endpointData, endpointDataError := bc.getEndpointData(tag)

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

func (bc Botocore) getEndpointData(tag github.RepoTag) (EndpointFile, error) {
	var partition EndpointFile

	blob, err := bc.Repo.GetBlobFromTag(tag.Name, "data/endpoints.json")

	if err != nil {
		return EndpointFile{}, err
	}

	if err := json.Unmarshal(blob, &partition); err != nil {
		return EndpointFile{}, err
	}

	return partition, nil
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
