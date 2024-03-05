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
}

// <---------------------------------------------------------------------------------------------------->

// New creates a new Conifg struct with the values from ./../configs/config.json
func New() (*Config, error) {
	newConfig := Config{}

	configMap, err := util.GetJSONData("./../configs/config.json")
	if err != nil {
		return &newConfig, fmt.Errorf("couldn't open config JSON file; %s", err.Error())
	}

	index := 0

	// populate the newConfig with properly formated strings
	for _, value := range configMap {
		newSlice := util.ConvertSliceInterface[string](value.([]interface{}))
		newSlice = util.FormatEntries(newSlice, true)
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
