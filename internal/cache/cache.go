// Package cache ...
package cache

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/SkillpTm/better-windows-search/internal/config"
	"github.com/SkillpTm/better-windows-search/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

var (
	EntrieFilesystem *Filesystem
)

// <---------------------------------------------------------------------------------------------------->

type Filesystem struct {
	mainDirs      map[string]map[int][][]interface{}
	secondaryDirs map[string]map[int][][]interface{}
}

// New returns a pointer to a Filesystem struct that has been filled up according to the config file
func New(mainDirPaths []string, secondaryDirPaths []string) *Filesystem {
	fs := Filesystem{mainDirs: make(map[string]map[int][][]interface{}), secondaryDirs: make(map[string]map[int][][]interface{})}

	fs.create(mainDirPaths, true)
	fs.create(secondaryDirPaths, false)

	return &fs
}

// create gets and sets the folders set to either mainDirPaths or secondaryDirPaths
func (fs *Filesystem) create(dirPaths []string, isMainDirs bool) {
	// if we aren't adding to the mainDirs add the excluded mainDirs directly to the queue
	if !isMainDirs {
		dirPaths = append(dirPaths, config.BWSConfig.ExcludeSubMainDirs...)
	}

	if len(dirPaths) < 1 {
		return
	}

	// 10000000 is the channel size, because we just need a ridiculously large channel to store all the paths until we traversed them
	var pathQueue = make(chan string, 10000000)

	for _, dir := range dirPaths {
		pathQueue <- dir
	}

	var wg sync.WaitGroup

	// 10000000 is the channel size, because we just need a ridiculously large channel to store all the results until we add them to the fs
	resultsChan := make(chan *[][]string, 10000000)

	for range config.BWSConfig.CPUThreads {
		wg.Add(1)
		go fs.traverse(pathQueue, isMainDirs, resultsChan, &wg)
	}

	wg.Wait()

	close(resultsChan)
	close(pathQueue)

	for result := range resultsChan {
		fs.add(result, isMainDirs)
	}
}

// walkDir walks through the pathQueue and adds all new and valid entries into the resultsChan
func (fs *Filesystem) traverse(pathQueue chan string, isMainDirs bool, resultsChan chan<- *[][]string, wg *sync.WaitGroup) {
	// when the queue is empty disolve the worker
	defer wg.Done()

	// loop over the queue until it's empty
	for currentDir := range pathQueue {
		newPaths := []string{}
		newEntries := [][]string{}

		currentEntries, err := os.ReadDir(currentDir)
		if err != nil {
			// an error here simply means we didn't have the permissions to read a dir, so we ignore it
			if len(pathQueue) < 1 {
				break
			}

			continue
		}

		for _, entry := range currentEntries {
			if entry.IsDir() {
				entryPath := util.FormatEntry(filepath.Join(currentDir, entry.Name()), true)

				// check if the current dir is an excluded name
				if util.SliceContains[string](config.BWSConfig.ExcludeDirsByName, util.FormatEntry(entry.Name(), true)) {
					continue
				}

				// check if the dir is excluded
				if util.SliceContains[string](config.BWSConfig.ExcludeDirs, entryPath) {
					continue
				}

				// check if we found a mainDirs folder while not mainDirs working with mainDirs
				if !isMainDirs && util.SliceContains[string](config.BWSConfig.Maindirs, entryPath) {
					continue
				}

				// check if the dir is in the excluded main dirs
				if isMainDirs && util.SliceContains[string](config.BWSConfig.ExcludeSubMainDirs, entryPath) {
					continue
				}

				newPaths = append(newPaths, entryPath)
				newEntries = append(newEntries, []string{entryPath, entry.Name(), "Folder"})
			} else {
				entryPath := util.FormatEntry(filepath.Join(currentDir, entry.Name()), false)

				fileExtension := filepath.Ext(entry.Name())
				if len(fileExtension) < 1 {
					fileExtension = "File"
				}
				newEntries = append(newEntries, []string{entryPath, entry.Name(), fileExtension})
			}
		}

		for _, path := range newPaths {
			pathQueue <- path
		}

		resultsChan <- &newEntries

		if len(pathQueue) < 1 {
			break
		}
	}
}

// add adds the newEntries to the fs
func (fs *Filesystem) add(newEntries *[][]string, isMainDirs bool) {
	var tempStorage map[string]map[int][][]interface{}

	if isMainDirs {
		tempStorage = fs.mainDirs
	} else {
		tempStorage = fs.secondaryDirs
	}

	for _, item := range *newEntries {
		itemPath := item[0]
		itemName := item[1]
		itemExtension := item[2]

		// trim file extensions from the name, if it has one
		if itemExtension != "File" && itemExtension != "Folder" {
			itemName = itemName[:len(itemName)-len(itemExtension)]
		}

		// check if the file type is already stored in the fs, if not add it in
		if _, ok := tempStorage[itemExtension]; !ok {
			tempStorage[itemExtension] = make(map[int][][]interface{})
		}

		// check if the file length is already stored for the file extension, if not add it in
		if _, ok := tempStorage[itemExtension][len(itemName)]; !ok {
			tempStorage[itemExtension][len(itemName)] = [][]interface{}{}
		}

		// add the file into the fs at its length with the format: [path, name, [encoded bytes]]
		tempStorage[itemExtension][len(itemName)] = append(tempStorage[itemExtension][len(itemName)], []interface{}{itemPath, itemName, Encode(itemName)})
	}

	if isMainDirs {
		fs.mainDirs = tempStorage
	} else {
		fs.secondaryDirs = tempStorage
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
