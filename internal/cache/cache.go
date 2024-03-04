// Package cache ...
package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

// <---------------------------------------------------------------------------------------------------->

var EntrieFilesystem = Filesystem{mainDirs: make(map[string]map[int]map[string][]interface{})}

// <---------------------------------------------------------------------------------------------------->

type Filesystem struct {
	mainDirs map[string]map[int]map[string][]interface{}
	// secondaryDirs map[string]map[int]map[string][]interface{}
}

func Create(rootPath string) {
	pathStack := []string{rootPath}

	for len(pathStack) > 0 {
		// pop dir from stack
		currentDir := pathStack[len(pathStack)-1]
		pathStack = pathStack[:len(pathStack)-1]

		currentEntries, err := os.ReadDir(currentDir)
		if err != nil {
			fmt.Printf("Error accessing directory %s: %v\n", currentDir, err)
			continue
		}

		tempSlice := [][]string{}

		for _, entry := range currentEntries {
			entryPath := filepath.Join(currentDir, entry.Name())

			if entry.IsDir() {
				pathStack = append(pathStack, entryPath)
			} else {
				tempSlice = append(tempSlice, []string{entryPath, entry.Name(), filepath.Ext(entry.Name())})
			}
		}
	}
}
