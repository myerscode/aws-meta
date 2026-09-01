package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveArchiveFile writes an archive data file beneath outputDir/pkg/data/archive.
func SaveArchiveFile(outputDir string, jsonData any, fileName string) error {
	return SaveData(jsonData, filepath.Join(outputDir, "pkg", "data", "archive", fileName))
}

// SaveManifestFile writes a manifest data file beneath outputDir/pkg/data/manifests.
func SaveManifestFile(outputDir string, jsonData any, fileName string) error {
	return SaveData(jsonData, filepath.Join(outputDir, "pkg", "data", "manifests", fileName))
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
