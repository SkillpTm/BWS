// Package cache ...
package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

		EntrieFilesystem.add(&tempSlice, true)
	}
}

func (fs *Filesystem) add(newSlice *[][]string, isMainDirs bool) {
	if !isMainDirs {
		return
	}

	for _, item := range *newSlice {
		itemPath := item[0]
		itemName := item[1]
		itemExtension := item[2]

		// trim file extensions from the name
		if len(itemExtension) > 0 {
			itemName = itemName[:len(itemName)-len(itemExtension)]
		} else {
			// if there is no file extension make the extension "File"
			itemExtension = "File"
		}

		// check if the file type is already stored in the fs, if not add it in
		if _, ok := fs.mainDirs[itemExtension]; !ok {
			fs.mainDirs[itemExtension] = make(map[int]map[string][]interface{})
		}

		// check if the file length is already stored for the file extension, if not add it in
		if _, ok := fs.mainDirs[itemExtension][len(itemName)]; !ok {
			fs.mainDirs[itemExtension][len(itemName)] = make(map[string][]interface{})
		}

		// add the file into the fs at its length with the path as a key and the name as it's value
		fs.mainDirs[itemExtension][len(itemName)][strings.ReplaceAll(itemPath, "\\", "/")] = []interface{}{itemName}
		//TODO add binary here
	}
}
