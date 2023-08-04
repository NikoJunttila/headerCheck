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
	flag.Parse()

	fmt.Println("checking files...")
	err = checkHeader(*projectPathFlagPtr, *forceFlagPtr, *yearFlagPtr, *authorFlagPtr)
	fmt.Println("All files checked")
	if err != nil {
		fmt.Printf("Error scanning project: %v\n", err)
		os.Exit(1)
	}
}
