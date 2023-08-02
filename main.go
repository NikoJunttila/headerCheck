package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: copyright_checker <project_path>")
		os.Exit(1)
	}

	projectPath := os.Args[1]

	err := filepath.Walk(projectPath, checkHeader)
	if err != nil {
		fmt.Printf("Error scanning project: %v\n", err)
		os.Exit(1)
	}
}
