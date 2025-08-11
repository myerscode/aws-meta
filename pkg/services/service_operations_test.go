package services

import (
	"fmt"
	"testing"
)

func TestServiceOperationsByServiceId(t *testing.T) {
	tests := []struct {
		name        string
		serviceId   string
		expectError bool
		expectEmpty bool
	}{
		{
			name:        "Valid service - ACM",
			serviceId:   "ACM",
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "Valid service - S3",
			serviceId:   "S3",
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "Valid service - EC2",
			serviceId:   "EC2",
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "Invalid service",
			serviceId:   "INVALID-SERVICE-123",
			expectError: true,
			expectEmpty: true,
		},
		{
			name:        "Empty service ID",
			serviceId:   "",
			expectError: true,
			expectEmpty: true,
		},
		{
			name:        "Case sensitive - lowercase",
			serviceId:   "acm",
			expectError: true,
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operations, err := ServiceOperationsByServiceId(tt.serviceId)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for service %s, but got none", tt.serviceId)
				}
				if len(operations) != 0 {
					t.Errorf("Expected empty operations list for invalid service, got %d operations", len(operations))
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for service %s: %v", tt.serviceId, err)
				}
				if tt.expectEmpty && len(operations) > 0 {
					t.Errorf("Expected empty operations list, got %d operations", len(operations))
				}
				if !tt.expectEmpty && len(operations) == 0 {
					t.Errorf("Expected non-empty operations list for service %s", tt.serviceId)
				}
			}
		})
	}
}

func TestServiceOperationsByServiceName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		expectError bool
		expectEmpty bool
	}{
		{
			name:        "Valid service - AWS Certificate Manager",
			serviceName: "AWS Certificate Manager",
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "Valid service - Amazon Simple Storage Service",
			serviceName: "Amazon Simple Storage Service",
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "Valid service - Amazon Elastic Compute Cloud",
			serviceName: "Amazon Elastic Compute Cloud",
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "Invalid service name",
			serviceName: "Invalid Service Name",
			expectError: true,
			expectEmpty: true,
		},
		{
			name:        "Empty service name",
			serviceName: "",
			expectError: true,
			expectEmpty: true,
		},
		{
			name:        "Case sensitive - lowercase",
			serviceName: "aws certificate manager",
			expectError: true,
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operations, err := ServiceOperationsByServiceName(tt.serviceName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for service %s, but got none", tt.serviceName)
				}
				if len(operations) != 0 {
					t.Errorf("Expected empty operations list for invalid service, got %d operations", len(operations))
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for service %s: %v", tt.serviceName, err)
				}
				if tt.expectEmpty && len(operations) > 0 {
					t.Errorf("Expected empty operations list, got %d operations", len(operations))
				}
				if !tt.expectEmpty && len(operations) == 0 {
					t.Errorf("Expected non-empty operations list for service %s", tt.serviceName)
				}
			}
		})
	}
}

// Test ServiceOperations alias function
func TestServiceOperations(t *testing.T) {
	// Test that the alias function works
	operations, err := ServiceOperations("ACM")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(operations) == 0 {
		t.Fatal("Expected operations for ACM")
	}

	// Test that it returns the same results as ServiceOperationsByServiceId
	operationsById, errById := ServiceOperationsByServiceId("ACM")
	if errById != nil {
		t.Fatalf("Unexpected error from ServiceOperationsByServiceId: %v", errById)
	}

	if len(operations) != len(operationsById) {
		t.Errorf("ServiceOperations returned different number of operations: %d vs %d", len(operations), len(operationsById))
	}

	for i, op := range operations {
		if i >= len(operationsById) || op != operationsById[i] {
			t.Errorf("Operations differ at index %d: %s vs %s", i, op, operationsById[i])
		}
	}
}

func TestServiceOperationsContentById(t *testing.T) {
	// Test that operations returned are valid strings for ServiceId
	operations, err := ServiceOperationsByServiceId("ACM")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(operations) == 0 {
		t.Fatal("Expected at least one operation for ACM")
	}

	// Check that all operations are non-empty strings
	for i, operation := range operations {
		if operation == "" {
			t.Errorf("Operation at index %d is empty string", i)
		}
	}

	// Check that operations are sorted (they should be based on the data generation)
	for i := 1; i < len(operations); i++ {
		if operations[i-1] > operations[i] {
			t.Errorf("Operations are not sorted: %s > %s", operations[i-1], operations[i])
		}
	}

	// Test specific known operations for ACM
	expectedOperations := []string{"AddTagsToCertificate", "DeleteCertificate", "DescribeCertificate"}
	for _, expected := range expectedOperations {
		found := false
		for _, operation := range operations {
			if operation == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected operation %s not found in ACM operations", expected)
		}
	}
}

func TestServiceOperationsContentByName(t *testing.T) {
	// Test that operations returned are valid strings for ServiceName
	operations, err := ServiceOperationsByServiceName("AWS Certificate Manager")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(operations) == 0 {
		t.Fatal("Expected at least one operation for AWS Certificate Manager")
	}

	// Check that all operations are non-empty strings
	for i, operation := range operations {
		if operation == "" {
			t.Errorf("Operation at index %d is empty string", i)
		}
	}

	// Check that operations are sorted
	for i := 1; i < len(operations); i++ {
		if operations[i-1] > operations[i] {
			t.Errorf("Operations are not sorted: %s > %s", operations[i-1], operations[i])
		}
	}
}

func TestServiceOperationsConsistency(t *testing.T) {
	// Test that both functions return the same operations for the same service
	testCases := []struct {
		serviceId   string
		serviceName string
	}{
		{"ACM", "AWS Certificate Manager"},
		{"S3", "Amazon Simple Storage Service"},
		{"EC2", "Amazon Elastic Compute Cloud"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%s", tc.serviceId, tc.serviceName), func(t *testing.T) {
			operationsById, errById := ServiceOperationsByServiceId(tc.serviceId)
			operationsByName, errByName := ServiceOperationsByServiceName(tc.serviceName)

			if errById != nil {
				t.Fatalf("Unexpected error for service ID %s: %v", tc.serviceId, errById)
			}
			if errByName != nil {
				t.Fatalf("Unexpected error for service name %s: %v", tc.serviceName, errByName)
			}

			if len(operationsById) != len(operationsByName) {
				t.Errorf("Different number of operations: ID=%d, Name=%d", len(operationsById), len(operationsByName))
			}

			for i, op := range operationsById {
				if i >= len(operationsByName) || op != operationsByName[i] {
					t.Errorf("Operations differ at index %d: %s vs %s", i, op, operationsByName[i])
				}
			}
		})
	}
}

func TestServiceOperationsMultipleServices(t *testing.T) {
	// Test multiple services to ensure they return different operations
	serviceIds := []string{"ACM", "S3", "EC2"}
	operationSets := make(map[string][]string)

	for _, serviceId := range serviceIds {
		operations, err := ServiceOperationsByServiceId(serviceId)
		if err != nil {
			t.Fatalf("Unexpected error for service %s: %v", serviceId, err)
		}
		if len(operations) == 0 {
			t.Fatalf("Expected operations for service %s", serviceId)
		}
		operationSets[serviceId] = operations
	}

	// Ensure different services have different operations (at least some)
	acmOps := operationSets["ACM"]
	s3Ops := operationSets["S3"]

	allSame := true
	if len(acmOps) != len(s3Ops) {
		allSame = false
	} else {
		for i, op := range acmOps {
			if op != s3Ops[i] {
				allSame = false
				break
			}
		}
	}

	if allSame {
		t.Error("ACM and S3 should have different operations")
	}
}
