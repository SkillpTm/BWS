// Package setup ...
package setup

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"

	"github.com/SkillpTm/BWS/internal/cache"
	"github.com/SkillpTm/BWS/internal/config"
)

// <---------------------------------------------------------------------------------------------------->

// Init sets up the BWSConfig and EntrieFilesystem
func Init() error {
	var err error

	config.BWSConfig, err = config.New()
	if err != nil {
		return fmt.Errorf("couldn't assign BWSConfig; %s", err.Error())
	}

	cache.EntrieFilesystem = cache.New(config.BWSConfig.Maindirs, config.BWSConfig.SecondaryDirs)

	return nil
}
