// Package util ...
package util

// <---------------------------------------------------------------------------------------------------->

// SliceContains checks if a slice includes a certain element
func SliceContains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}

	return false
}
