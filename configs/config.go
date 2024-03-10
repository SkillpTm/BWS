// Package defaultconfig ...
package defaultconfig

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/json"
	"fmt"
)

// <---------------------------------------------------------------------------------------------------->

var rawJSONData = []byte(`
{
    "cpuThreads": 4,
    "mainDirs": [
        "C:/Users/<USERNAME>/"
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
}`)

// <---------------------------------------------------------------------------------------------------->

// Get will provide a map with the config JSON data
func Get() (map[string]interface{}, error) {
	var jsonData map[string]interface{}

	err := json.Unmarshal(rawJSONData, &jsonData)
	if err != nil {
		return jsonData, fmt.Errorf("couldn't unmarshal raw JSON data; %s", err.Error())
	}

	return jsonData, nil
}
