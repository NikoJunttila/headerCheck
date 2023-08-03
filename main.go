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
	flag.Parse()
	fmt.Println("checking files...")
	err = checkHeader(*projectPathFlagPtr, *forceFlagPtr)
	fmt.Println("All files checked")
	if err != nil {
		fmt.Printf("Error scanning project: %v\n", err)
		os.Exit(1)
	}
}
