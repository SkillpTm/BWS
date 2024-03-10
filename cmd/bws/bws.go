// Package bws ...
package bws

// <---------------------------------------------------------------------------------------------------->

import (
	"time"

	"github.com/skillptm/bws/internal/cache"
	"github.com/skillptm/bws/internal/config"
	"github.com/skillptm/bws/internal/search"
	"github.com/skillptm/bws/internal/setup"
)

// <---------------------------------------------------------------------------------------------------->

func init() {
	setup.Init()

	go updateCache()
}

func updateCache() {
	// create tickers for how often the FileSystem components are supossed to update
	mainDirsTicker := time.NewTicker(3 * time.Minute)
	defer mainDirsTicker.Stop()

	secondaryDirsTicker := time.NewTicker(30 * time.Minute)
	defer secondaryDirsTicker.Stop()

	for {
		select {
		case <-mainDirsTicker.C:
			// check if the Filesystem is Updateable, if so update it, otherwise wait for the next cycle
			if cache.EntrieFilesystem.Updateable {
				cache.EntrieFilesystem.Update(config.BWSConfig.Maindirs, true)
			}
		case <-secondaryDirsTicker.C:
			// check if the Filesystem is Updateable, if so update it, otherwise wait for the next cycle
			if cache.EntrieFilesystem.Updateable {
				cache.EntrieFilesystem.Update(config.BWSConfig.SecondaryDirs, false)
			}
		}
	}
}

// Search takes in any substring that you want to search for through all filenames. You may add any amount if file extensions as well.
// The extendedSearch flag dictates, if we search through the SecondaryDirs from the ./configs/config.json file.
// To change the folders included/excluded in the search edit the ./configs/config.json file.
// This function is a wrapper around the search and rank functions and returns the ranked search results at the end.
func Search(searchString string, fileExtensions []string, extendedSearch bool) *[]string {
	// wait for the FileSystem to be readable
	for {
		if cache.EntrieFilesystem.Readable {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	// make it so while we search we can't update the FileSystem
	cache.EntrieFilesystem.Updateable = false
	defer func() {
		cache.EntrieFilesystem.Updateable = true
	}()

	// get the filepaths and names
	results, pattern := search.Start(search.NewSearchString(searchString, fileExtensions), extendedSearch)

	// rank and sort the files
	return search.Rank(results, pattern)
}
