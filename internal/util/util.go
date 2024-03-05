// Package util ...
package util

// <---------------------------------------------------------------------------------------------------->

// <---------------------------------------------------------------------------------------------------->

// ConvertSliceInterface converts a []interface{} to a slice of any type
func ConvertSliceInterface[T any](sliceInput []interface{}) []T {
	newSlice := []T{}

	for _, item := range sliceInput {
		newSlice = append(newSlice, item.(T))
	}

	return newSlice
}
