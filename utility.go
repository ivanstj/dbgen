package main

import (
	"encoding/json"
	"os"
)

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// parse json file
func parseDatatype() (result map[string]string) {
	// Open our jsonFile

	json.Unmarshal([]byte(postgresToGolang), &result)
	return
}
