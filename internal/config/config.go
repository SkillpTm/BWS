// Package config ...
package config

import (
	"fmt"

	"github.com/SkillpTm/better-windows-search/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

var BWSConfig *Config

// <---------------------------------------------------------------------------------------------------->

type Config struct {
	Maindirs           []string
	ExcludeSubMainDirs []string
	ExcludeDirs        []string
	ExcludeDirsByName  []string
}

// <---------------------------------------------------------------------------------------------------->

// New creates a new Conifg struct with the values from ./configs/config.json
func New() (*Config, error) {
	newConfig := Config{}

	configMap, err := util.GetJSONData("./configs/config.json")
	if err != nil {
		return &newConfig, fmt.Errorf("couldn't open config JSON file; %s", err.Error())
	}

	index := 0

	// populate the newConfig with properly formated paths
	for _, value := range configMap {
		newSlice := util.ConvertSliceInterface[string](value.([]interface{}))

		// on index 3 we  only want the slice of strings and not format them, since they aren't paths
		if index == 3 {
			newConfig.ExcludeDirsByName = newSlice
			continue
		}

		for index, element := range newSlice {
			newSlice[index] = util.FormatEntry(element, true)
		}
		newSlice, err = util.InsertUsername(newSlice)
		if err != nil {
			return &newConfig, fmt.Errorf("couldn't replace '<USERNAME>'; %s", err.Error())
		}

		switch index {
		case 0:
			newConfig.Maindirs = newSlice
		case 1:
			newConfig.ExcludeSubMainDirs = newSlice
		case 2:
			newConfig.ExcludeDirs = newSlice
		}

		index++
	}

	return &newConfig, nil
}
