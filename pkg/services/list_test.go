package services

import (
	"testing"
)

func TestServiceManifest(t *testing.T) {
	manifest, err := serviceManifest()

	if err != nil {
		t.Errorf("error reading manifest: %v", err)
	}

	if manifest == nil {
		t.Errorf("Manifest not loaded: %v", err)
	}
}

func TestAllServiceNames(t *testing.T) {
	serviceNames := AllServiceNames()

	if len(serviceNames) == 0 {
		t.Errorf("AllServiceNames() returned no service names")
	}
	t.Logf("AllServiceNames() returned %d service names", len(serviceNames))
}
