// Package cache handles everything that has to do with the generation of the cache for the Search function.
package cache

// <---------------------------------------------------------------------------------------------------->

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/skillptm/ssl/pkg/sslslices"

	"github.com/skillptm/bws/internal/config"
	"github.com/skillptm/bws/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

var (
	EntrieFilesystem *Filesystem = &Filesystem{SetupProperly: false}
)

// <---------------------------------------------------------------------------------------------------->

type Filesystem struct {
	MainDirs      map[string]map[int][][]interface{}
	SecondaryDirs map[string]map[int][][]interface{}

	SetupProperly bool
	Readable      bool
	Updateable    bool
}

// New returns a pointer to a Filesystem struct that has been filled up according to the config file
func New(mainDirPaths []string, secondaryDirPaths []string) *Filesystem {
	fs := Filesystem{
		MainDirs:      make(map[string]map[int][][]interface{}),
		SecondaryDirs: make(map[string]map[int][][]interface{}),
		SetupProperly: false,
		Readable:      false,
		Updateable:    false,
	}

	fs.Update(mainDirPaths, true)
	fs.Update(secondaryDirPaths, false)

	return &fs
}

// Update gets and sets the folders set to either mainDirPaths or secondaryDirPaths
func (fs *Filesystem) Update(dirPaths []string, isMainDirs bool) {
	fs.Updateable = false

	// if we aren't adding to the MainDirs add the excluded MainDirs directly to the queue
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
	breakChan := make(chan bool, config.BWSConfig.CPUThreads*config.BWSConfig.CPUThreads)

	for range config.BWSConfig.CPUThreads {
		wg.Add(1)
		go fs.traverse(isMainDirs, pathQueue, resultsChan, breakChan, &wg)
	}

	wg.Wait()

	close(resultsChan)
	close(pathQueue)
	close(breakChan)

	// make a tempFS, in case there is something on the current FS, so there'll always be a cache for the search to grab
	tempFS := Filesystem{MainDirs: make(map[string]map[int][][]interface{}), SecondaryDirs: make(map[string]map[int][][]interface{})}

	for result := range resultsChan {
		tempFS.add(result, isMainDirs)
	}

	// while overwritting set the state to unreadable until we're done
	fs.Readable = false

	// overwrite the old cache with the new cache
	if isMainDirs {
		fs.MainDirs = tempFS.MainDirs
	} else {
		fs.SecondaryDirs = tempFS.SecondaryDirs
	}

	fs.Readable = true
	fs.Updateable = true
}

// walkDir walks through the pathQueue and adds all new and valid entries into the resultsChan
func (fs *Filesystem) traverse(isMainDirs bool, pathQueue chan string, resultsChan chan<- *[][]string, breakChan chan bool, wg *sync.WaitGroup) {
	// when the queue is empty disolve the worker
	defer wg.Done()

	for {
		select {
		// loop over the queue until it's empty
		case currentDir := <-pathQueue:
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
					if sslslices.Contains[string](config.BWSConfig.ExcludeDirsByName, util.FormatEntry(entry.Name(), true)) {
						continue
					}

					// check if the dir is excluded
					if sslslices.Contains[string](config.BWSConfig.ExcludeDirs, entryPath) {
						continue
					}

					// check if we found a MainDirs folder while not MainDirs working with MainDirs
					if !isMainDirs && sslslices.Contains[string](config.BWSConfig.MainDirs, entryPath) {
						continue
					}

					// check if the dir is in the excluded main dirs
					if isMainDirs && sslslices.Contains[string](config.BWSConfig.ExcludeSubMainDirs, entryPath) {
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
				for range config.BWSConfig.CPUThreads {
					breakChan <- true
				}
				return
			}

		case <-breakChan:
			return
		}
	}
}

// add adds the newEntries to the fs
func (fs *Filesystem) add(newEntries *[][]string, isMainDirs bool) {
	var tempStorage map[string]map[int][][]interface{}

	if isMainDirs {
		tempStorage = fs.MainDirs
	} else {
		tempStorage = fs.SecondaryDirs
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
		tempStorage[itemExtension][len(itemName)] = append(tempStorage[itemExtension][len(itemName)], []interface{}{itemPath, strings.ToLower(itemName), Encode(itemName)})
	}

	if isMainDirs {
		fs.MainDirs = tempStorage
	} else {
		fs.SecondaryDirs = tempStorage
	}
}
