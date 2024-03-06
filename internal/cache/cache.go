// Package cache ...
package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SkillpTm/better-windows-search/internal/config"
	"github.com/SkillpTm/better-windows-search/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

var EntrieFilesystem *Filesystem

// <---------------------------------------------------------------------------------------------------->

type Filesystem struct {
	mainDirs      map[string]map[int][][]interface{}
	secondaryDirs map[string]map[int][][]interface{}
}

// New returns a pointer to a Filesystem struct that has been filled up according to the config file
func New(dirPaths *[]string, isMainDirs bool) *Filesystem {
	fs := Filesystem{mainDirs: make(map[string]map[int][][]interface{})}

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
			if entry.IsDir() {
				entryPath := util.FormatEntry(filepath.Join(currentDir, entry.Name()), true)

				// check if the current dir is an excluded name
				if util.SliceContains[string](config.BWSConfig.ExcludeDirsByName, entry.Name()) {
					continue
				}

				// check if the dir is excluded
				if util.SliceContains[string](config.BWSConfig.ExcludeDirs, entryPath) {
					continue
				}

				// check if the dir is in the excluded main dirs
				if isMainDirs && util.SliceContains[string](config.BWSConfig.ExcludeSubMainDirs, entryPath) {
					continue
				}

				pathStack = append(pathStack, entryPath)
				tempSlice = append(tempSlice, []string{entryPath, entry.Name(), "Folder"})
			} else {
				entryPath := util.FormatEntry(filepath.Join(currentDir, entry.Name()), false)

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
			fs.mainDirs[itemExtension] = make(map[int][][]interface{})
		}

		// check if the file length is already stored for the file extension, if not add it in
		if _, ok := fs.mainDirs[itemExtension][len(itemName)]; !ok {
			fs.mainDirs[itemExtension][len(itemName)] = [][]interface{}{}
		}

		// add the file into the fs at its length with the path as a key and the name as it's value
		fs.mainDirs[itemExtension][len(itemName)] = append(fs.mainDirs[itemExtension][len(itemName)], []interface{}{util.FormatEntry(itemPath, false), itemName})
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
