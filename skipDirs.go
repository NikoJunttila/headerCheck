package main

// add more when needed
var foldersToSkip = []string{
	"node_modules",
	".venv",
	"build",
  "vendor",
}
var filesToSkip = []string{}

func shouldSkipDirOrFile(name string, isDir bool) bool {
	if isDir {
		for _, folder := range foldersToSkip {
			if name == folder {
				return true
			}
		}
	} else {
		for _, file := range filesToSkip {
			if name == file {
				return true
			}
		}
	}
	return false
}
