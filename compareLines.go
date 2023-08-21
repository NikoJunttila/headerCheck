/****************************************************************
*
* File   : compareLines.go
* Author : Niko Junttila <niko.junttila2@centria.fi>
*          NikoJunttila <89527972+NikoJunttila@users.noreply.github.com>
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
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func showDifferences(newLines []string, oldLines []string) {
	diff := len(newLines) - len(oldLines)
	insertLen := len(oldLines) - 13
	var result []string
	if diff > 0 {
		newString := "*"

		insertIndex := 5 + insertLen
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
}
