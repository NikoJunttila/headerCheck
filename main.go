package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	defaultProjectPath, err := os.Getwd()
	projectPathFlagPtr := flag.String("src", defaultProjectPath, "a string")
	forceFlagPtr := flag.Bool("force", false, "a bool")
	authorFlagPtr := flag.String("author", "default", "a string")
	yearFlagPtr := flag.String("year", "2023", "a string")
	vControlPtr := flag.String("control", "git", "a string")
	flag.Parse()

	fmt.Println("checking files...")
	if *vControlPtr == "git" {
		readIgnore()
		err = gitCheckHeader(*projectPathFlagPtr, *forceFlagPtr, *yearFlagPtr, *authorFlagPtr)
	} else if *vControlPtr == "mer" {
		fmt.Println("work in progress. try again later :(")
	} else {
		fmt.Printf("Unexpected version control. check if there is typo in `%v` \n", *vControlPtr)
		fmt.Println(`currently accepted commands are --control="git" and --control="mer" default=git`)
		os.Exit(1)
	}

	fmt.Println("All files checked")
	if err != nil {
		fmt.Printf("Error scanning project: %v\n", err)
		os.Exit(1)
	}
}
