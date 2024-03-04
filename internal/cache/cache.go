// Package cache ...
package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// <---------------------------------------------------------------------------------------------------->

var EntrieFilesystem Filesystem

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
				// do check for is excluded at all, if so just don't add to stack
				pathStack = append(pathStack, entryPath)
				tempSlice = append(tempSlice, []string{entryPath, entry.Name(), "Folder"})
			} else {
				fileExtension := filepath.Ext(entry.Name())
				if len(fileExtension) < 1 {
					fileExtension = "File"
				}
				tempSlice = append(tempSlice, []string{entryPath, entry.Name(), fileExtension})
			}
		}

		fs.add(&tempSlice, isMainDirs)
	}
}

func (fs *Filesystem) add(newFiles *[][]string, isMainDirs bool) {
	if !isMainDirs {
		return
	}

	for _, item := range *newFiles {
		itemPath := item[0]
		itemName := item[1]
		itemExtension := item[2]

		// trim file extensions from the name, if it has one
		if itemExtension != "File" && itemExtension != "Folder" {
			itemName = itemName[:len(itemName)-len(itemExtension)]
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
