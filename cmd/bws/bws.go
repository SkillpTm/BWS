package main

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
	"time"

	"github.com/SkillpTm/BWS/internal/search"
	"github.com/SkillpTm/BWS/internal/setup"
)

// <---------------------------------------------------------------------------------------------------->

func main() {
	startTime := time.Now()

	err := setup.Init()
	fmt.Println(err)

	elapsedTime := time.Since(startTime)

	fmt.Println("Make Map:", elapsedTime)

	time.Sleep(15 * time.Second)

	startTime = time.Now()

	search.Start("Haribo", []string{"Folder"}, true)

	elapsedTime = time.Since(startTime)

	fmt.Println("Search Time:", elapsedTime)

	for {
		time.Sleep(1 * time.Minute)
	}
}
