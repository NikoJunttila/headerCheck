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

func showBlockDifferences(newLines []string, oldLines []string){
    oldLines, newLines = removeCommonElements(oldLines, newLines)
    for i := 0; i < len(oldLines); i++ {
        color.Red("Old line: - %s" ,oldLines[i])
    }
    for j := 0; j < len(newLines); j++{
      color.Green("new Line: + %s", newLines[j])
    }
}

func showDifferences(newLines []string, oldLines []string, lines int, authInsert int) error {
	diff := len(newLines) - len(oldLines)
	insertLen := len(oldLines) - lines
	var result []string
	if diff > 0 {
		newString := " *"
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
    
    showOldLines, showNewLines := removeCommonElements(oldLines, newLines)
    
    for i := 0; i < len(showOldLines); i++ {
        color.Red("Old line: - %s" ,showOldLines[i])
    }
    for i := 0; i < len(showNewLines); i++{
      color.Green("new Lines: + %s", showNewLines[i])
    }
    return nil
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
  
    var trimmed1 []string
    var trimmed2 []string

    
    for _, str := range slice2{
      trimmed2 = append(trimmed2, strings.TrimSpace(str))
    }
    for _, str := range slice1{
      trimmed1 = append(trimmed1, strings.TrimSpace(str))
    }
    for _, element := range slice1 {
      element = strings.TrimSpace(element)
        if !contains(element, trimmed2) {
            resultSlice1 = append(resultSlice1, element)
        }
    }
    for _, element := range slice2 {
      element = strings.TrimSpace(element)
        if !contains(element, trimmed1) {
            resultSlice2 = append(resultSlice2, element)
        }
    }
    return resultSlice1, resultSlice2
}
