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
	CPUThreads         int
	Maindirs           []string
	ExcludeSubMainDirs []string
	SecondaryDirs      []string
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

	newConfig.CPUThreads = int(configMap["cpuThreads"].(float64))
	delete(configMap, "cpuThreads")

	// populate the newConfig with properly formated paths
	for key, value := range configMap {
		newSlice := util.ConvertSliceInterface[string](value.([]interface{}))

		for index, element := range newSlice {
			newSlice[index] = util.FormatEntry(element, true)
		}
		newSlice, err = util.InsertUsername(newSlice)
		if err != nil {
			return &newConfig, fmt.Errorf("couldn't replace '<USERNAME>'; %s", err.Error())
		}

		switch key {
		case "mainDirs":
			newConfig.Maindirs = newSlice
		case "excludeSubMainDirs":
			newConfig.ExcludeSubMainDirs = newSlice
		case "secondaryDirs":
			newConfig.SecondaryDirs = newSlice
		case "excludeDirs":
			newConfig.ExcludeDirs = newSlice
		case "excludeDirsByName":
			newConfig.ExcludeDirsByName = newSlice
		}
	}

	return &newConfig, nil
}
