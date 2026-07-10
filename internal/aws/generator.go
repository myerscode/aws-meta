package aws

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

func SaveArchiveFile(jsonData any, fileName string) error {
	metaDataFile := fmt.Sprintf("pkg/data/archive/%s", fileName)
	return SaveData(jsonData, metaDataFile)
}

func SaveManifestFile(jsonData any, fileName string) error {
	metaDataFile := fmt.Sprintf("pkg/data/manifests/%s", fileName)
	return SaveData(jsonData, metaDataFile)
}

func SaveData(jsonData any, fileName string) error {
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(jsonData, "", " ")
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	if err := os.WriteFile(fileName, data, 0644); err != nil {
		return fmt.Errorf("error writing file %s: %w", fileName, err)
	}

	return nil
}
