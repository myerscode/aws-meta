package util

import (
	"fmt"
	"reflect"
	"sort"
)

// Sort returns a new slice with the input strings sorted alphabetically.
func Sort(input []string) []string {
	sorted := make([]string, len(input))
	copy(sorted, input)
	sort.Strings(sorted)
	return sorted
}

// SortByField sorts a slice of structs by a named string field.
func SortByField(slice any, fieldName string) error {
	v := reflect.ValueOf(slice)

	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("expected pointer to slice, got %T", slice)
	}

	s := v.Elem()

	sort.Slice(s.Interface(), func(i, j int) bool {
		a := s.Index(i).FieldByName(fieldName)
		b := s.Index(j).FieldByName(fieldName)

		if !a.IsValid() || !b.IsValid() {
			return false
		}

		if a.Kind() != reflect.String || b.Kind() != reflect.String {
			return false
		}

		return a.String() < b.String()
	})

	return nil
}
