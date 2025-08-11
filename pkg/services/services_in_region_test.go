package services

import (
	"testing"
)

func TestServicesInRegion(t *testing.T) {
	tests := []struct {
		name        string
		regionName  string
		expectError bool
		expectEmpty bool
	}{
		{
			name:        "Valid region - us-east-1",
			regionName:  "us-east-1",
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "Valid region - eu-west-1",
			regionName:  "eu-west-1",
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "Invalid region",
			regionName:  "invalid-region-123",
			expectError: true,
			expectEmpty: true,
		},
		{
			name:        "Empty region name",
			regionName:  "",
			expectError: true,
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			services, err := ServicesInRegion(tt.regionName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for region %s, but got none", tt.regionName)
				}
				if len(services) != 0 {
					t.Errorf("Expected empty services list for invalid region, got %d services", len(services))
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for region %s: %v", tt.regionName, err)
				}
				if tt.expectEmpty && len(services) > 0 {
					t.Errorf("Expected empty services list, got %d services", len(services))
				}
				if !tt.expectEmpty && len(services) == 0 {
					t.Errorf("Expected non-empty services list for region %s", tt.regionName)
				}
			}
		})
	}
}

func TestServicesInRegionContent(t *testing.T) {

	services, err := ServicesInRegion("us-east-1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(services) == 0 {
		t.Fatal("Expected at least one service in us-east-1")
	}

	// Check that all services are non-empty strings
	for i, service := range services {
		if service == "" {
			t.Errorf("Service at index %d is empty string", i)
		}
	}

	// Check that services are sorted (they should be based on the data generation)
	for i := 1; i < len(services); i++ {
		if services[i-1] > services[i] {
			t.Errorf("Services are not sorted: %s > %s", services[i-1], services[i])
		}
	}
}
