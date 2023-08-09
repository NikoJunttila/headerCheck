package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	defaultProjectPath, err := os.Getwd()
	projectPathFlagPtr := flag.String("src", defaultProjectPath, "a string")
	forceFlagPtr := flag.Bool("force", false, "a bool")
	authorFlagPtr := flag.String("author", "default", "a string")
	yearFlagPtr := flag.String("year", "2023", "a string")
	/* 	vControlPtr := flag.String("control", "git", "a string")*/
	flag.Parse()

	fmt.Println("checking files...")
	dotGitfile := filepath.Join(defaultProjectPath, ".git")
	_, err = os.Stat(dotGitfile)
	if err == nil {
		fmt.Println("found .git")
		readIgnore(".gitignore")
		err = gitCheckHeader(*projectPathFlagPtr, *forceFlagPtr, *yearFlagPtr, *authorFlagPtr)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("defaulted to mercurial")
		readIgnore(".hgignore") //idk idk fr fr
		err = mercuCheckHeader(*projectPathFlagPtr, *forceFlagPtr, *yearFlagPtr, *authorFlagPtr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("All files checked")
	if err != nil {
		fmt.Printf("Error scanning project: %v\n", err)
		os.Exit(1)
	}
}
