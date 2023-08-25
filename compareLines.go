/****************************************************************
 *
 *  File   : compareLines.go
 *  Author : NikoJunttila <89527972+NikoJunttila@users.noreply.github.com>
 *           Niko Junttila <niko.junttila2@centria.fi>
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
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func showDifferences(newLines []string, oldLines []string, lines int, authInsert int) error {
	diff := len(newLines) - len(oldLines)
	insertLen := len(oldLines) - lines
	var result []string
	if diff > 0 {
		newString := "*"
		//append * to old header so checking lines stays even at index 5 + authors???
		// modify hardcoded 5 to suit new template if changing author position
    insertIndex := authInsert + insertLen
		if insertIndex < 1 {
			return errors.New("header messed up")
		}
		before := oldLines[:insertIndex]
		after := oldLines[insertIndex:]
		result = append(result, before...)
		result = append(result, newString)
		for i := 1; i < diff; i++ {
			result = append(result, newString)
		}
		result = append(result, after...)
	} else if diff < 0 {
    //in case there are less authors in new header
    color.Yellow("Removed authors here")
    oldLines, newLines := removeCommonElements(oldLines, newLines)
    for i := 0; i < len(oldLines); i++ {
        color.Red("Old line: - %s" ,oldLines[i])
    }
    for i := 0; i < len(newLines); i++{
      color.Green("new Lines: + %s", newLines[i])
    }
	} else{
		result = oldLines
  }
	//compare all lines of both templates and show difference
	for i := 0; i < len(result) && i < len(newLines); i++ {
		resultLine := strings.TrimSpace(result[i])
		newLine := strings.TrimSpace(newLines[i])
		resultLine = strings.ToLower(resultLine)
		newLine = strings.ToLower(newLine)
		if resultLine != newLine {
			color.Red("Line: %d - %s", i+1, result[i])
			color.Green("Line: %d + %s", i+1, newLines[i])
			fmt.Println()
		}
	}
	return nil
}

func removeCommonElements(slice1, slice2 []string) ([]string, []string) {

    var resultSlice1, resultSlice2 []string

    // Create a map to store the elements of slice2 for faster lookup
    elementsInSlice2 := make(map[string]bool)
    for _, element := range slice2 {
        elementsInSlice2[element] = true
    }

    // Iterate through slice1 and check if each element is in slice2
    for _, element := range slice1 {
        if _, exists := elementsInSlice2[element]; !exists {
            resultSlice1 = append(resultSlice1, element)
        }
    }

    // Iterate through slice2 and check if each element is in slice1
    for _, element := range slice2 {
        if _, exists := elementsInSlice2[element]; !exists {
            resultSlice2 = append(resultSlice2, element)
        }
    }

    return resultSlice1, resultSlice2
}
