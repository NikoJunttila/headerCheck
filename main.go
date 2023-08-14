/****************************************************************
*
* File   : main.go
* Author : NikoJunttila <89527972+NikoJunttila@users.noreply.github.com>
*
*
* Copyright (C) 2023 Centria University of Applied Sciences.
* All rights reserved.
*
* Unauthorized copying of this file, via any medium is strictly
* prohibited.
*
****************************************************************/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func main() {
	defaultProjectPath, err := os.Getwd()
	forceFlagPtr := flag.Bool("force", false, "a bool")
	authorFlagPtr := flag.String("author", "default", "default author if no repo histories")
	yearFlagPtr := flag.String("year", "2023", "default year if no repo histories")
	ignoreFolderFlagPtr := flag.String("ignore", "", "ignore folder")
	flag.Parse()
	if len(*ignoreFolderFlagPtr) > 1 {
		foldersToSkip = append(foldersToSkip, *ignoreFolderFlagPtr)
	}
	fmt.Println("checking files...")
	dotGitfile := filepath.Join(defaultProjectPath, ".git")
	_, err = os.Stat(dotGitfile)
	if err == nil {
		color.Green("found .git")
		readIgnore(".gitignore")
		err = gitCheckHeader(defaultProjectPath, *forceFlagPtr, *yearFlagPtr, *authorFlagPtr)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("defaulted to mercurial")
		readIgnore(".hgignore") // idk idk fr fr
		err = mercuCheckHeader(defaultProjectPath, *forceFlagPtr, *yearFlagPtr, *authorFlagPtr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("All files checked")
	if err != nil {
		color.Red("Error scanning project: %v\n", err)
		os.Exit(1)
	}
}
