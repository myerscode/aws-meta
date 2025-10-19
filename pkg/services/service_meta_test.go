package services

import (
	"testing"
)

func TestServiceMeta(t *testing.T) {
	services, err := ServiceMeta()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(services) == 0 {
		t.Fatal("Expected at least one service")
	}

	// Check that all services have required fields
	for serviceId, service := range services {
		if service.ServiceId == "" {
			t.Errorf("Service %s has empty ServiceId", serviceId)
		}
		if service.ServiceFullName == "" {
			t.Errorf("Service %s has empty ServiceFullName", service.ServiceId)
		}
		if len(service.Operations) == 0 {
			t.Errorf("Service %s has no operations", service.ServiceId)
		}

		// Check that all operations are non-empty strings
		for j, operation := range service.Operations {
			if operation == "" {
				t.Errorf("Service %s, operation at index %d is empty", service.ServiceId, j)
			}
		}
	}
}

func TestServiceMetaContent(t *testing.T) {
	services, err := ServiceMeta()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check for some known services
	expectedServices := []string{"ACM", "S3", "EC2", "Lambda", "RDS"}

	for _, expected := range expectedServices {
		if _, found := services[expected]; !found {
			t.Errorf("Expected service %s not found", expected)
		}
	}

	// Verify that service IDs match map keys
	for serviceId, service := range services {
		if service.ServiceId != serviceId {
			t.Errorf("Service ID mismatch: map key %s != service.ServiceId %s", serviceId, service.ServiceId)
		}
	}
}

func TestServiceMetaSpecificData(t *testing.T) {
	services, err := ServiceMeta()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test S3 service
	s3Service, s3Found := services["S3"]
	if !s3Found {
		t.Fatal("S3 service not found")
	}
	if s3Service.ServiceFullName != "Amazon Simple Storage Service" {
		t.Errorf("Expected S3 service full name to be 'Amazon Simple Storage Service', got %s", s3Service.ServiceFullName)
	}
	if len(s3Service.Operations) == 0 {
		t.Error("Expected S3 service to have operations")
	}

	// Test ACM service
	acmService, acmFound := services["ACM"]
	if !acmFound {
		t.Fatal("ACM service not found")
	}
	if acmService.ServiceFullName != "AWS Certificate Manager" {
		t.Errorf("Expected ACM service full name to be 'AWS Certificate Manager', got %s", acmService.ServiceFullName)
	}
	if len(acmService.Operations) == 0 {
		t.Error("Expected ACM service to have operations")
	}

	// Test that S3 has more operations than ACM (generally true)
	if len(s3Service.Operations) <= len(acmService.Operations) {
		t.Errorf("Expected S3 to have more operations than ACM: %d vs %d",
			len(s3Service.Operations), len(acmService.Operations))
	}
}

func TestServiceMetaOperations(t *testing.T) {
	services, err := ServiceMeta()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test S3 service operations
	s3Service, s3Found := services["S3"]
	if !s3Found {
		t.Fatal("S3 service not found")
	}

	// Check for some known S3 operations
	expectedOperations := []string{"GetObject", "PutObject", "ListBuckets", "CreateBucket"}
	foundOperations := make(map[string]bool)

	for _, operation := range s3Service.Operations {
		foundOperations[operation] = true
	}

	for _, expected := range expectedOperations {
		if !foundOperations[expected] {
			t.Errorf("Expected S3 operation %s not found", expected)
		}
	}

	// Check that operations within a service are sorted
	for i := 1; i < len(s3Service.Operations); i++ {
		if s3Service.Operations[i-1] > s3Service.Operations[i] {
			t.Errorf("Operations in S3 service are not sorted: %s > %s",
				s3Service.Operations[i-1], s3Service.Operations[i])
		}
	}
}

func TestServiceMetaConsistency(t *testing.T) {
	// Test that ServiceMeta returns the same data as other service functions
	services, err := ServiceMeta()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	allServiceNames := AllServiceNames()

	// Check that the count matches
	if len(services) != len(allServiceNames) {
		t.Errorf("ServiceMeta returned %d services, but AllServiceNames returned %d", len(services), len(allServiceNames))
	}

	// Check that all service names from AllServiceNames are present in ServiceMeta
	serviceNamesFromMeta := make(map[string]bool)
	for _, service := range services {
		serviceNamesFromMeta[service.ServiceFullName] = true
	}

	for _, serviceName := range allServiceNames {
		if !serviceNamesFromMeta[serviceName] {
			t.Errorf("Service name %s from AllServiceNames not found in ServiceMeta", serviceName)
		}
	}
}

func TestServiceMetaForRegion(t *testing.T) {
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
			services, err := ServiceMetaForRegion(tt.regionName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for region %s, but got none", tt.regionName)
				}
				if len(services) != 0 {
					t.Errorf("Expected empty services map for invalid region, got %d services", len(services))
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for region %s: %v", tt.regionName, err)
				}
				if tt.expectEmpty && len(services) > 0 {
					t.Errorf("Expected empty services map, got %d services", len(services))
				}
				if !tt.expectEmpty && len(services) == 0 {
					t.Errorf("Expected non-empty services map for region %s", tt.regionName)
				}
			}
		})
	}
}

func TestServiceMetaForRegionContent(t *testing.T) {
	services, err := ServiceMetaForRegion("us-east-1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(services) == 0 {
		t.Fatal("Expected at least one service in us-east-1")
	}

	// Check that all services have required fields
	for serviceId, service := range services {
		if service.ServiceId == "" {
			t.Errorf("Service %s has empty ServiceId", serviceId)
		}
		if service.ServiceFullName == "" {
			t.Errorf("Service %s has empty ServiceFullName", service.ServiceId)
		}
		if len(service.Operations) == 0 {
			t.Errorf("Service %s has no operations", service.ServiceId)
		}

		// Verify that service ID matches map key
		if service.ServiceId != serviceId {
			t.Errorf("Service ID mismatch: map key %s != service.ServiceId %s", serviceId, service.ServiceId)
		}
	}
}

func TestServiceMetaForRegionConsistency(t *testing.T) {
	// Test that ServiceMetaForRegion returns consistent data with ServicesInRegion
	regionName := "us-east-1"

	serviceEndpoints, err := ServicesInRegion(regionName)
	if err != nil {
		t.Fatalf("Unexpected error from ServicesInRegion: %v", err)
	}

	serviceMeta, err := ServiceMetaForRegion(regionName)
	if err != nil {
		t.Fatalf("Unexpected error from ServiceMetaForRegion: %v", err)
	}

	// ServiceMetaForRegion should return services that are available in the region
	// Note: There can be more services than endpoint prefixes because some endpoint prefixes
	// map to multiple services (e.g., "rds" maps to RDS, Neptune, and DocDB)
	if len(serviceMeta) == 0 {
		t.Error("ServiceMetaForRegion should return at least some services for us-east-1")
	}

	// Create a map of endpoint prefixes available in the region
	regionEndpointSet := make(map[string]bool)
	for _, endpoint := range serviceEndpoints {
		regionEndpointSet[endpoint] = true
	}

	// Verify that all services returned by ServiceMetaForRegion have endpoints available in the region
	for serviceId, service := range serviceMeta {
		if !regionEndpointSet[service.EndpointPrefix] {
			t.Errorf("Service %s (endpoint: %s) returned by ServiceMetaForRegion is not available in region %s",
				serviceId, service.EndpointPrefix, regionName)
		}
	}

	// Test that some known services that should be in us-east-1 are present
	knownServices := []string{"s3", "ec2", "lambda", "rds"}
	foundServices := 0
	for _, knownEndpoint := range knownServices {
		for _, service := range serviceMeta {
			if service.EndpointPrefix == knownEndpoint {
				foundServices++
				break
			}
		}
	}

	if foundServices < len(knownServices) {
		t.Errorf("Expected to find %d known services in us-east-1, but only found %d",
			len(knownServices), foundServices)
	}
}

func TestServiceMetaForRegionSpecificData(t *testing.T) {
	services, err := ServiceMetaForRegion("us-east-1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test that some known services are present in us-east-1
	expectedServices := []string{"S3", "EC2", "Lambda"}

	for _, expectedServiceId := range expectedServices {
		if service, found := services[expectedServiceId]; found {
			// Verify the service has the expected structure
			if service.ServiceId != expectedServiceId {
				t.Errorf("Expected service ID %s, got %s", expectedServiceId, service.ServiceId)
			}
			if service.ServiceFullName == "" {
				t.Errorf("Service %s has empty ServiceFullName", expectedServiceId)
			}
			if len(service.Operations) == 0 {
				t.Errorf("Service %s has no operations", expectedServiceId)
			}
		} else {
			t.Errorf("Expected service %s not found in us-east-1", expectedServiceId)
		}
	}
}

func TestServiceMetaForRegionSubset(t *testing.T) {
	// Test that ServiceMetaForRegion returns a subset of ServiceMeta
	allServices, err := ServiceMeta()
	if err != nil {
		t.Fatalf("Unexpected error from ServiceMeta: %v", err)
	}

	regionServices, err := ServiceMetaForRegion("us-east-1")
	if err != nil {
		t.Fatalf("Unexpected error from ServiceMetaForRegion: %v", err)
	}

	// Check that region services is a subset of all services
	if len(regionServices) > len(allServices) {
		t.Errorf("ServiceMetaForRegion returned more services (%d) than ServiceMeta (%d)",
			len(regionServices), len(allServices))
	}

	// Check that all services in region are also in the complete service list
	for serviceId, regionService := range regionServices {
		if allService, found := allServices[serviceId]; found {
			// Verify they are the same service
			if regionService.ServiceFullName != allService.ServiceFullName {
				t.Errorf("Service %s has different ServiceFullName: region=%s, all=%s",
					serviceId, regionService.ServiceFullName, allService.ServiceFullName)
			}
			if len(regionService.Operations) != len(allService.Operations) {
				t.Errorf("Service %s has different operation count: region=%d, all=%d",
					serviceId, len(regionService.Operations), len(allService.Operations))
			}
		} else {
			t.Errorf("Service %s from ServiceMetaForRegion not found in ServiceMeta", serviceId)
		}
	}
}
