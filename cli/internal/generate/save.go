package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

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
