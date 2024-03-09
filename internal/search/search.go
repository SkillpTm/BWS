// Package search ...
package search

// <---------------------------------------------------------------------------------------------------->

import (
	"strings"

	"github.com/SkillpTm/better-windows-search/internal/cache"
	"github.com/SkillpTm/better-windows-search/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

// SearchString holds all the data releated to the searchString input, so we only have to calculate them once
type SearchString struct {
	encoded    [8]byte
	extensions []string
	length     int
	name       string
}

// Start wrapps around the search and rank functions and returns the ranked search results at the end
func Start(searchString string, fileExtensions []string, extendedSearch bool) *[]string {
	results, pattern := search(newSearchString(searchString, fileExtensions), extendedSearch)

	return rank(results, pattern)
}

// newSearchString returns a pointer to a SearchString struct based on the string input
func newSearchString(searchString string, fileExtensions []string) *SearchString {
	// make sure all extensions begin with a period, unless it's a "File" or a "Folder"
	for index, element := range fileExtensions {
		if string(element[0]) != "." && element != "File" && element != "Folder" {
			fileExtensions[index] = "." + element
			continue
		}

		// ensure "File"/"Folder" have the right case
		if string(element[0]) != "f" {
			fileExtensions[index] = "F" + element[1:]
		}
	}

	return &SearchString{
		encoded:    cache.Encode(searchString),
		extensions: fileExtensions,
		length:     len(searchString),
		name:       strings.ToLower(searchString),
	}
}

// search wraps around the walkFS function and returns all the results from the MainDirs and SecondaryDirs
func search(pattern *SearchString, extendedSearch bool) (*[][]string, *SearchString) {
	output := [][]string{}

	// check the MainDirs for the search string
	output = append(output, *pattern.searchFS(&cache.EntrieFilesystem.MainDirs)...)

	// check the SecondaryDirs for the search string
	if extendedSearch {
		output = append(output, *pattern.searchFS(&cache.EntrieFilesystem.SecondaryDirs)...)
	}

	return &output, pattern
}

// searchFS searches one of the provided FileSystem maps, while skiping files for wrong extensions and ecoded values
func (searchString *SearchString) searchFS(dirs *map[string]map[int][][]interface{}) *[][]string {
	output := [][]string{}

	// loop over the extensions
	for extenion, lengthMaps := range *dirs {
		// check if extensions were provided and if so, if the current extension is a provided one
		if len(searchString.extensions) > 0 && !util.SliceContains[string](searchString.extensions, extenion) {
			continue
		}

		// loop over the filename lengths
		for length, fileSlices := range lengthMaps {
			// check if the filename is longer than the searchString
			if length < searchString.length {
				continue
			}

			// loop over the actual files
			for _, file := range fileSlices {
				// check if all required letters are inside the filename
				if !cache.CompareBytes(searchString.encoded, file[2].([8]byte)) {
					continue
				}

				// do a substring search over the filename
				if !strings.Contains(file[1].(string), searchString.name) {
					continue
				}

				// if the searchString is inside the filename add it's path and name to the output
				output = append(output, []string{file[0].(string), file[1].(string)})
			}
		}
	}

	return &output
}
