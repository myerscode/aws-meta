package data

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/myerscode/aws-meta/internal/aws"

	"github.com/Masterminds/semver/v3"
)

//go:embed manifests/*
var manifestsFS embed.FS

var manifestDirectoryName = "manifests"

// jsonData, err := getArchiveFile(`botocore.partitions\.(.+)\.json`)
func getArchiveFile(pattern string) ([]byte, error) {

	files, err := manifestsFS.ReadDir(manifestDirectoryName)

	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	re := regexp.MustCompile(pattern)

	var latestVersion *semver.Version

	var latestFile string

	for _, file := range files {
		matches := re.FindStringSubmatch(file.Name())
		if len(matches) < 2 {
			continue
		}

		version, err := semver.NewVersion(matches[1])
		if err != nil {
			continue
		}

		if latestVersion == nil || version.GreaterThan(latestVersion) {
			latestVersion = version
			latestFile = file.Name()
		}
	}

	if latestFile == "" {
		return nil, fmt.Errorf("no valid JSON files found")
	}

	content, err := manifestsFS.ReadFile(filepath.Join(manifestDirectoryName, latestFile))
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", latestFile, err)
	}

	return content, nil
}

func getManifestFile(file string) ([]byte, error) {
	// Read the file directly from the manifestDirectoryName location
	filePath := filepath.Join(manifestDirectoryName, file)
	content, err := manifestsFS.ReadFile(filePath)

	if err != nil {
		return nil, fmt.Errorf("error reading manifest file %s: %w", filePath, err)
	}

	return content, nil
}

func GetLatestParitionArchiveFile() (aws.PartitionSchemas, error) {
	jsonData, err := getArchiveFile(`botocore.partitions\.(.+)\.json`)

	if err != nil {
		return nil, err
	}

	var schema aws.PartitionSchemas

	if err := json.Unmarshal(jsonData, &schema); err != nil {
		return nil, err
	}

	return schema, nil
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
