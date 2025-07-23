package aws

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/myerscode/aws-meta/internal/util"
)

type ServiceSchemas map[string]ServiceSchema
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

func TouchResultsFile(path string) {

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		util.PrintErrorAndExit(err)
	}

	file, err := os.Create(path)

	if err != nil {
		util.PrintErrorAndExit(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			util.PrintErrorAndExit(err)
		}
	}(file)
}

func SaveArchiveFile(jsonData interface{}, fileName string) error {
	metaDataFile := fmt.Sprintf("pkg/data/archive/%s", fileName)

	return SaveData(jsonData, metaDataFile)
}

func SaveManifestFile(jsonData interface{}, fileName string) error {
	metaDataFile := fmt.Sprintf("pkg/data/manifests/%s", fileName)

	return SaveData(jsonData, metaDataFile)
}

func SaveData(jsonData interface{}, fileName string) error {

	TouchResultsFile(fileName)

	data, _ := json.MarshalIndent(jsonData, "", " ")

	if err := os.WriteFile(fileName, data, 0644); err != nil {
		util.PrintErrorAndExit(err)
	}

	return nil
}
