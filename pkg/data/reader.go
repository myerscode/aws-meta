package data

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"

	"github.com/myerscode/aws-meta/internal/aws"
)

//go:embed manifests/*
var manifestsFS embed.FS

var manifestDirectoryName = "manifests"

func getManifestFile(file string) ([]byte, error) {
	filePath := path.Join(manifestDirectoryName, file)
	content, err := manifestsFS.ReadFile(filePath)

	if err != nil {
		return nil, fmt.Errorf("error reading manifest file %s: %w", filePath, err)
	}

	return content, nil
}

func PartitionManifest() (aws.PartitionSchemas, error) {
	jsonData, err := getManifestFile(`botocore.partitions.json`)

	if err != nil {
		return nil, err
	}

	var testData aws.PartitionSchemas

	if err := json.Unmarshal(jsonData, &testData); err != nil {
		return nil, err
	}

	return testData, nil
}

func RegionsManifest() (aws.RegionSchemas, error) {
	jsonData, err := getManifestFile(`botocore.regions.json`)

	if err != nil {
		return nil, err
	}

	var testData aws.RegionSchemas

	if err := json.Unmarshal(jsonData, &testData); err != nil {
		return nil, err
	}

	return testData, nil
}

func ServiceManifest() (aws.ServiceSchemas, error) {
	jsonData, err := getManifestFile(`botocore.services.json`)

	if err != nil {
		return nil, err
	}

	var testData aws.ServiceSchemas

	if err := json.Unmarshal(jsonData, &testData); err != nil {
		return nil, err
	}

	return testData, nil
}
