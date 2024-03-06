// Package cache ...
package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SkillpTm/better-windows-search/internal/config"
)

// <---------------------------------------------------------------------------------------------------->

var EntrieFilesystem *Filesystem

// <---------------------------------------------------------------------------------------------------->

type Filesystem struct {
	mainDirs map[string]map[int]map[string][]interface{}
	// secondaryDirs map[string]map[int]map[string][]interface{}
}

func New(dirPaths *[]string, isMainDirs bool) *Filesystem {
	fs := Filesystem{mainDirs: make(map[string]map[int]map[string][]interface{})}

	fs.createDirs(dirPaths, isMainDirs)
	// fs.createDirs(rootPath, false)

	return &fs
}

func (fs *Filesystem) createDirs(dirPaths *[]string, isMainDirs bool) {
	pathStack := *dirPaths

	if !isMainDirs {
		return
		// exclude from main dirs add to secondary stack
	}

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
				// check if the dir is excluded
				if checkDirExcluded(entryPath, config.BWSConfig.ExcludeDirs) {
					continue
				}

				// check if the dir is in the excluded main dirs
				if isMainDirs && checkDirExcluded(entryPath, config.BWSConfig.ExcludeSubMainDirs) {
					continue
				}

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

// returns true if the a dir is in the provided slice
func checkDirExcluded(dir string, excludedDirs []string) bool {
	for _, excludedDir := range excludedDirs {
		if strings.ReplaceAll(dir, "\\", "/") == excludedDir[:len(excludedDir)-1] {
			return true
		}
	}

	return false
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

// WriteToFile saves mainDirs and secondaryDirs into individual jsons
func (fs *Filesystem) WriteToFile() error {
	err := util.WriteToJSON("./cache/mainDirs.json", fs.mainDirs)
	if err != nil {
		return fmt.Errorf("couldn't write mainDirs to JSON; %s", err.Error())
	}

	err = util.WriteToJSON("./cache/secondaryDirs.json", fs.secondaryDirs)
	if err != nil {
		return fmt.Errorf("couldn't write secondaryDirs to JSON; %s", err.Error())
	}

	return nil
}
