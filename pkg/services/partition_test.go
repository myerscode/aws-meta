package services

import (
	"testing"
)

func TestIsolatedPartitions(t *testing.T) {
	partitions, err := IsolatedPartitions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(partitions) == 0 {
		t.Fatal("Expected at least one isolated partition")
	}

	// Check that all returned partitions contain "iso"
	for _, partition := range partitions {
		if !contains(partition, "iso") {
			t.Errorf("Partition %s should contain 'iso' to be considered isolated", partition)
		}
	}

	// Check for known isolated partitions
	expectedPartitions := []string{"aws-iso", "aws-iso-b", "aws-iso-e", "aws-iso-f"}
	for _, expected := range expectedPartitions {
		if !containsString(partitions, expected) {
			t.Errorf("Expected isolated partition %s not found", expected)
		}
	}

	// Ensure no commercial or sovereign partitions are included
	excludedPartitions := []string{"aws", "aws-cn", "aws-us-gov", "aws-eusc"}
	for _, excluded := range excludedPartitions {
		if containsString(partitions, excluded) {
			t.Errorf("Isolated partitions should not include %s", excluded)
		}
	}
}

func TestSovereignPartitions(t *testing.T) {
	partitions, err := SovereignPartitions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(partitions) == 0 {
		t.Fatal("Expected at least one sovereign partition")
	}

	// Check for known sovereign partitions
	expectedPartitions := []string{"aws-cn", "aws-us-gov", "aws-eusc"}
	for _, expected := range expectedPartitions {
		if !containsString(partitions, expected) {
			t.Errorf("Expected sovereign partition %s not found", expected)
		}
	}

	// Ensure no commercial partitions are included
	if containsString(partitions, "aws") {
		t.Error("Sovereign partitions should not include commercial partition 'aws'")
	}

	// Ensure no isolated partitions are included
	for _, partition := range partitions {
		if contains(partition, "iso") {
			t.Errorf("Sovereign partitions should not include isolated partition %s", partition)
		}
	}
}

func TestCommercialPartitions(t *testing.T) {
	partitions, err := CommercialPartitions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(partitions) == 0 {
		t.Fatal("Expected at least one commercial partition")
	}

	// Should contain exactly the "aws" partition
	if len(partitions) != 1 {
		t.Errorf("Expected exactly 1 commercial partition, got %d", len(partitions))
	}

	if partitions[0] != "aws" {
		t.Errorf("Expected commercial partition to be 'aws', got %s", partitions[0])
	}

	// Ensure no sovereign or isolated partitions are included
	excludedPartitions := []string{"aws-cn", "aws-us-gov", "aws-eusc", "aws-iso", "aws-iso-b", "aws-iso-e", "aws-iso-f"}
	for _, excluded := range excludedPartitions {
		if containsString(partitions, excluded) {
			t.Errorf("Commercial partitions should not include %s", excluded)
		}
	}
}

func TestPartitionCategorization(t *testing.T) {
	// Test that all partitions are categorized into exactly one category
	allPartitions := AllPartitionNames()

	isolated, err := IsolatedPartitions()
	if err != nil {
		t.Fatalf("Error getting isolated partitions: %v", err)
	}

	sovereign, err := SovereignPartitions()
	if err != nil {
		t.Fatalf("Error getting sovereign partitions: %v", err)
	}

	commercial, err := CommercialPartitions()
	if err != nil {
		t.Fatalf("Error getting commercial partitions: %v", err)
	}

	// Count total categorized partitions
	totalCategorized := len(isolated) + len(sovereign) + len(commercial)

	if totalCategorized != len(allPartitions) {
		t.Errorf("Total categorized partitions (%d) doesn't match all partitions (%d)", totalCategorized, len(allPartitions))
	}

	// Ensure no partition appears in multiple categories
	allCategorized := make(map[string]int)

	for _, p := range isolated {
		allCategorized[p]++
	}
	for _, p := range sovereign {
		allCategorized[p]++
	}
	for _, p := range commercial {
		allCategorized[p]++
	}

	for partition, count := range allCategorized {
		if count != 1 {
			t.Errorf("Partition %s appears in %d categories, should appear in exactly 1", partition, count)
		}
	}

	// Ensure every partition from AllPartitionNames is categorized
	for _, partition := range allPartitions {
		if _, found := allCategorized[partition]; !found {
			t.Errorf("Partition %s from AllPartitionNames is not categorized", partition)
		}
	}
}

func TestPartitionContent(t *testing.T) {
	// Test that partitions are non-empty strings
	testCases := []struct {
		name string
		fn   func() ([]string, error)
	}{
		{"IsolatedPartitions", IsolatedPartitions},
		{"SovereignPartitions", SovereignPartitions},
		{"CommercialPartitions", CommercialPartitions},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			partitions, err := tc.fn()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			for i, partition := range partitions {
				if partition == "" {
					t.Errorf("Partition at index %d is empty string", i)
				}
			}

			// Check that partitions are sorted (they should be based on the data generation)
			for i := 1; i < len(partitions); i++ {
				if partitions[i-1] > partitions[i] {
					t.Errorf("Partitions are not sorted: %s > %s", partitions[i-1], partitions[i])
				}
			}
		})
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || (len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
