# Better Windows Search [![Go Report Card](https://goreportcard.com/badge/github.com/skillptm/bws)](https://goreportcard.com/report/github.com/skillptm/bws)

BWS is a Go module that generates a cache, which allows for way fast search times than the refula windows search. It also will always find the file, if it exists and will never *bing* your search.

## Build Struture:

The module creates a cache before the first search (which means the first search may take a while, all searches afterwards will be very fast though).

The cache is seperated in 2 different maps: MainDirs and SecondaryDirs. MainDirs get always searched and more frequently updated, while SecondaryDirs can only be searched with the extenedSearch flag and get updated more rarely.

There is a default config that you can update with the set functions in ./pkg/options. The default config looks liké this (it's not actually in a JSON):
```jsonc
{
	"cpuThreads": "1/4 of threads (int)", // this is set to the rounded up integer of 1/4 of your CPU threads
	"mainDirs": [
		"C:/Users/<USERNAME>/" // all instances of <USERNAME> get automatically repleased by the module, you can insert it like this too
    ],
	"excludeSubMainDirs": [
		"C:/Users/<USERNAME>/AppData/Roaming"
    ],
	"secondaryDirs": [
		"C:/"
    ],
	"excludeDirs": [
		"C:/Windows/",
		"C:/$Recycle.Bin/",
		"C:/Users/<USERNAME>/AppData/Local",
		"C:/Users/<USERNAME>/AppData/LocalLow"
    ],
	"excludeDirsByName": [
		".git",
		"bin",
		"node_modules",
		"steamapps"
    ]
}
```

## Usage:

The only functions in this module are:
- bws.Search: Used for a regular search
- bws.GoSearchWithBreak: The same as Search, just with the option of ending it early (which will cause you receiving an empty result)
- bws.ForceUpdateCache: No matter the circumstances updates the cache.
- The Set functions inside pkg/options used to change the modules config.

### Example:

```go
package main

import (
	"fmt"

	"github.com/skillptm/bws"
	"github.com/skillptm/bws/pkg/options"
)

func main() {
	options.SetCPUThreads(8)

	catResults := bws.Search("cat video", []string{"mp4", "mkv"}, false)

	for _, result := range catResults {
		fmt.Println(result)
	}

	dogResults := bws.Search("dog pictures", []string{"Folder"}, true)

	for _, result := range dogResults {
		fmt.Println(result)
	}

	var dragonResults []string
	breakChan := make(chan bool, 1)

	go func() {
		dragonResults = bws.GoSearchWithBreak("dragon audio", []string{".mp3"}, true, breakChan)
	}()

	breakChan <- true // now the goroutine has stopped and put an empty slice on dragonResults

	for _, result := range dragonResults {
		fmt.Println(result) // this line won't be executed
	}
}
```