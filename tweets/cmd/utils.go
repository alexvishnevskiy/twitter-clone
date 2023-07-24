package main

import (
	"os"
	"path/filepath"
)

func getStoragePath() string {
	// Get the current directory
	currentDir, _ := os.Getwd()
	// Traverse two levels up
	twoLevelsUp := filepath.Join(currentDir, "..", "..")
	// Get the absolute path
	absTwoLevelsUp, _ := filepath.Abs(twoLevelsUp)
	// Append storage to the path
	storagePath := filepath.Join(absTwoLevelsUp, "storage")
	return storagePath
}
