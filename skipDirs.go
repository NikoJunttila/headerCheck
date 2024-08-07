/****************************************************************
 *
 *  File   : skipDirs.go
 *  Author : NikoJunttila <89527972+NikoJunttila@users.noreply.github.com>
 *
 *  Copyright (C) 2023-2024 Centria University of Applied Sciences.
 *  All rights reserved.
 *
 *  Unauthorized copying of this file, via any medium is strictly
 *  prohibited.
 *
 ****************************************************************/

package main

// add more when needed
var foldersToSkip = []string{
	"node_modules",
	".venv",
	"build",
  "vendor",
  "venv",
}

var filesToSkip = []string{}

func shouldSkipDirOrFile(name string, isDir bool) bool {
  filesToSkip = append(filesToSkip, foldersToSkip...)
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
