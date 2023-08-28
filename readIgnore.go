/****************************************************************
 *
 *  File   : readIgnore.go
 *  Author : NikoJunttila <89527972+NikoJunttila@users.noreply.github.com>
 *
 *  Copyright (C) 2023 Centria University of Applied Sciences.
 *  All rights reserved.
 *
 *  Unauthorized copying of this file, via any medium is strictly
 *  prohibited.
 *
 ****************************************************************/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readIgnore(gitOrHg string) {
	gitignorePath := gitOrHg
	file, err := os.Open(gitignorePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var directories []string
	var files []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := strings.TrimSpace(scanner.Text())

		if len(entry) == 0 || strings.HasPrefix(entry, "#") {
			continue
		}
		if strings.HasPrefix(entry, "/") {
			entry = entry[1:]
		}
		info, err := os.Stat(entry)
		if err != nil {
			continue
		}

		if info.IsDir() {
			directories = append(directories, entry)
		} else {
			files = append(files, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return
	}
	foldersToSkip = append(foldersToSkip, directories...)
	filesToSkip = append(filesToSkip, files...)
}
