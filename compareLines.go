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

func showDifferences(newLines []string, oldLines []string, lines int) error {
	diff := len(newLines) - len(oldLines)
	insertLen := len(oldLines) - lines
	var result []string
	if diff > 0 {
		newString := "*"
		//append * to old header so checking lines stays even at index 5 + authors???
		// modify hardcoded 5 to suit new template if changing author position
    insertIndex := 5 + insertLen
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
	} else {
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
