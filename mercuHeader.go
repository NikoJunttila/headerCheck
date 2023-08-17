/****************************************************************
*
* File   : mercuHeader.go
* Author : NikoJunttila <89527972+NikoJunttila@users.noreply.github.com>
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
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func mercuCheckHeader(rootDir string, force bool, yearFlag string, authorFlag string) error {
	err := filepath.WalkDir(rootDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if shouldSkipDirOrFile(info.Name(), info.IsDir()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		suffix := filepath.Ext(path)
		var templateContent string
		switch {
		case contains(suffix, defaultSuffix):
			templateContent = templates[0].Header
		case contains(suffix, pySuffix):
			templateContent = templates[1].Header
		case contains(suffix, htmlSuffix):
			templateContent = templates[2].Header
		default:
			return nil
		}


    relative, err := filepath.Rel(rootDir, path) 
    fmt.Println(relative)
    if err != nil {
      fmt.Println("error with relative path")
      return err
    }
		// filename := filepath.Base(path)
		filenameModded := "'" + relative + "'"



		var trimmedYearRange string
		cmd := exec.Command("hg", "log", "--template", "{date|shortdate}\n", "-r", "reverse(ancestors(file("+filenameModded+")))")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error running 'hg log' command for file %s: %v\nOutput: %s\n", path, err, output)
			trimmedYearRange = yearFlag
		} else {
			commitDates := strings.Fields(string(output))
			var years []string
			for _, date := range commitDates {
				years = append(years, date[:4])
			}
			years = deduplicateAndSort(years)
			yearRange := formatYearRange(years)
			trimmedYearRange = strings.ReplaceAll(yearRange, `"`, "")
		}

		var trimmedAuthorList string
		var authors []string
		cmd2 := exec.Command("hg", "log", "--template", "{author|person} <{author|email}>\n", path)
		output2, err := cmd2.Output()
		if err != nil {
			trimmedAuthorList = authorFlag
		} else {
			authors = deduplicateAndSort(strings.Split(strings.TrimSpace(string(output2)), "\n"))
			authorList := strings.Join(authors, "\n*          ")
			trimmedAuthorList = strings.ReplaceAll(authorList, `"`, "")
		}

		templateContent = strings.ReplaceAll(templateContent, "{YEARS}", trimmedYearRange)
		templateContent = strings.ReplaceAll(templateContent, "{AUTHOR}", trimmedAuthorList)
		templateContent = strings.ReplaceAll(templateContent, "{FILENAME}", filepath.Base(path))

		existingContent, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			return err
		}
		existingHeader := ""
		maxLines := 18 + len(authors)
		headerLinesSplit := strings.Split(string(existingContent), "\n")
		if len(headerLinesSplit) < maxLines {
			maxLines = len(headerLinesSplit)
		}
		switch {
		case contains(suffix, defaultSuffix):
			for i := 0; i < maxLines; i++ {
				line := headerLinesSplit[i]
				if strings.Contains(
					line,
					"**********************************************************/",
				) {
					headerStartIndex := strings.Index(
						string(existingContent),
						"**********************************************************/",
					)
					if headerStartIndex != -1 {
						existingHeader = string(
							existingContent[:headerStartIndex+len("**********************************************************/")],
						)
						existingContent = existingContent[headerStartIndex+len("**********************************************************/"):]
					}
				}
			}
		case contains(suffix, pySuffix):
			for i := 0; i < maxLines; i++ {
				line := headerLinesSplit[i]
				if strings.Contains(
					line,
					`********************************************************"""`,
				) {
					headerStartIndex := strings.Index(
						string(existingContent),
						`********************************************************"""`,
					)
					if headerStartIndex != -1 {
						existingHeader = string(
							existingContent[:headerStartIndex+len(`********************************************************"""`)],
						)
						existingContent = existingContent[headerStartIndex+len(`********************************************************"""`):]
					}
				}
			}
		case contains(suffix, htmlSuffix):
			for i := 0; i < maxLines; i++ {
				line := headerLinesSplit[i]
				if strings.Contains(
					line,
					`---------------------------------------------------------->`,
				) {
					headerStartIndex := strings.Index(
						string(existingContent),
						`---------------------------------------------------------->`,
					)
					if headerStartIndex != -1 {
						existingHeader = string(
							existingContent[:headerStartIndex+len(`---------------------------------------------------------->`)],
						)
						existingContent = existingContent[headerStartIndex+len(`---------------------------------------------------------->`):]
					}
				}
			}
		default:
			fmt.Println("error no suffix found")
			return nil
		}
  	cleanedHeader := cleanString(existingHeader) 
		cleanedtemplateContent := cleanString(templateContent)
		if cleanedHeader == cleanedtemplateContent {
			fmt.Printf("File %s is good \n", path)
			return nil
		}
		if !force && len(existingHeader) < 10 {
			color.Red("No header found: %s \n", path)
			return nil
		}
		if !force {
			color.Red("file %s needs fix \n  \n", path)
			oldLines := strings.Split(existingHeader, "\n")
			newLines := strings.Split(templateContent, "\n")
      showDifferences(newLines, oldLines)
			return nil
		}
		// Combine the new header with the existing content
		newContent := templateContent + "\n" + string(existingContent)

		// Write the updated content back to the file
		err = os.WriteFile(path, []byte(newContent), 0644)
		if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", path, err)
			return err
		}

		color.Green("Copyright header fixed for file: %s\n", path)
		return nil
	})
	return err
}
