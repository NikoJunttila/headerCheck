package main

import (
	"os"
)

// add more when needed
var foldersToSkip = []string{
	"node_modules",
	".venv",
	"build",
}

func shouldSkipDir(d os.DirEntry) bool {
	if d.IsDir() {
		dirName := d.Name()
		for _, folder := range foldersToSkip {
			if dirName == folder {
				return true
			}
		}
	}
	return false
}
