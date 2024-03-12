// Package bws containes the main Search function and start up logic of bws.
package bws

// <---------------------------------------------------------------------------------------------------->

import (
	"time"

	"github.com/skillptm/bws/internal/cache"
	"github.com/skillptm/bws/internal/config"
	"github.com/skillptm/bws/internal/search"
)

// <---------------------------------------------------------------------------------------------------->

// init executes as soon as bws is imported into another project and launches a goroutine of updateCache
func init() {
	go updateCache()
}

// updateCache updates with teh use of tickers both the MainDirs and the SecondaryDirs
func updateCache() {
	// create tickers for how often the FileSystem components are supossed to update
	mainDirsTicker := time.NewTicker(3 * time.Minute)
	defer mainDirsTicker.Stop()

	secondaryDirsTicker := time.NewTicker(30 * time.Minute)
	defer secondaryDirsTicker.Stop()

	for {
		select {
		case <-mainDirsTicker.C:
			if !cache.EntrieFilesystem.SetupProperly {
				continue
			}

			// check if the Filesystem is Updateable, if so update it, otherwise wait for the next cycle
			if cache.EntrieFilesystem.Updateable {
				cache.EntrieFilesystem.Update(config.BWSConfig.MainDirs, true)
			}
		case <-secondaryDirsTicker.C:
			if !cache.EntrieFilesystem.SetupProperly {
				continue
			}

			// check if the Filesystem is Updateable, if so update it, otherwise wait for the next cycle
			if cache.EntrieFilesystem.Updateable {
				cache.EntrieFilesystem.Update(config.BWSConfig.SecondaryDirs, false)
			}
		}
	}
}

/*
Search takes in any substring that you want to search for through all filenames. You may add any amount of file extensions as well.
The extendedSearch flag dictates, if we search through the SecondaryDirs.
To change the folders included/excluded in the search use the pkg/options set functions.
This function is a wrapper around the search and rank functions and returns the ranked search results at the end.

On it's first execution the function will take longer, as it needs to generate the cache first.
*/
func Search(searchString string, fileExtensions []string, extendedSearch bool) []string {
	// check if the FileSystem is setup properly, if not reset it by regenerating it
	if !cache.EntrieFilesystem.SetupProperly {
		cache.EntrieFilesystem = cache.New(config.BWSConfig.MainDirs, config.BWSConfig.SecondaryDirs)
	}

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
	return *search.Rank(results, pattern)
}
