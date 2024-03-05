// Package util ...
package util

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

// <---------------------------------------------------------------------------------------------------->

// ConvertSliceInterface converts a []interface{} to a slice of any type
func ConvertSliceInterface[T any](sliceInput []interface{}) []T {
	newSlice := []T{}

	for _, item := range sliceInput {
		newSlice = append(newSlice, item.(T))
	}

	return newSlice
}

// InsertUsername replaces <USERNAME> in any path to the actual username of the current user
func InsertUsername(pathInputs []string) ([]string, error) {
	// Get the current user for their name
	currentUser, err := user.Current()
	if err != nil {
		return pathInputs, fmt.Errorf("couldn't get current user for username; %s", err.Error())
	}

	username := strings.Split(currentUser.Username, "\\")[1]

	for index, path := range pathInputs {
		pathInputs[index] = strings.ReplaceAll(path, "<USERNAME>", username)
	}

	return pathInputs, nil
}

// FormatEntries replaces "\" with "/" and for folders adds a "/" at the end, if there isn't one already
// This can't be used with a mix of files and folders
func FormatEntries(entries []string, isFolder bool) []string {
	for index, entry := range entries {
		entries[index] = strings.ReplaceAll(entry, "\\", "/")

		// if the inputs are just entries skip them
		if !isFolder {
			continue
		}

		if string(entry[len(entry)-1]) != "/" {
			entries[index] = entries[index] + "/"
		}
	}

	return entries
}

// GetJSONData will provide a map with JSON data of the provided file
func GetJSONData(filePath string) (map[string]interface{}, error) {
	// define return var early
	var jsonData map[string]interface{}
	var returnErr error = nil

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return jsonData, fmt.Errorf("couldn't open JSON file (%s); %s", filePath, err.Error())
	}

	// defer close the file with error handling
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			returnErr = fmt.Errorf("couldn't close JSON file (%s); %s", filePath, jsonFile.Close())
		}
	}()

	defer jsonFile.Close()

	rawJSONData, err := io.ReadAll(jsonFile)
	if err != nil {
		return jsonData, fmt.Errorf("couldn't read JSON file (%s); %s", filePath, err.Error())
	}

	err = json.Unmarshal(rawJSONData, &jsonData)
	if err != nil {
		return jsonData, fmt.Errorf("couldn't unmarshal raw JSON data; %s", err.Error())
	}

	return jsonData, returnErr
}
