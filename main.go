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
	authorFlagPtr := flag.String("author", "default", "a string")
	yearFlagPtr := flag.String("year", "2023", "a string")
	flag.Parse()

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
		readIgnore(".hgignore") //idk idk fr fr
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
