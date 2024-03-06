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
func ConvertSliceInterface[T comparable](sliceInput []interface{}) []T {
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

// FormatEntry replaces "\" with "/" and for folders adds a "/" at the end, if there isn't one already
func FormatEntry(entry string, isFolder bool) string {
	entry = strings.ReplaceAll(entry, "\\", "/")

	// if the inputs are just entries skip them
	if !isFolder {
		return entry
	}

	// check if the last char is a "/" if not append it
	if string(entry[len(entry)-1]) != "/" {
		entry += "/"
	}

	return entry
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

// WriteToJSON will take a map with JSON formated data and the file path and write to that file
func WriteToJSON(filePath string, inputData interface{}) error {
	// define return var early
	var returnErr error = nil

	jsonFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("couldn't open JSON file (%s); %s", filePath, err.Error())
	}

	// defer close the file with error handling
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			returnErr = fmt.Errorf("couldn't close JSON file (%s); %s", filePath, jsonFile.Close())
		}
	}()

	jsonData, err := json.MarshalIndent(inputData, "", "	")
	if err != nil {
		return fmt.Errorf("couldn't marshal JSON data; %s", err.Error())
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("couldn't write JSON data to JSON file (%s); %s", filePath, err.Error())
	}

	return returnErr
}
