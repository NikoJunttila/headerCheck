/****************************************************************
 *
 *  File   : mercuHeader.go
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
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
  "sort"

	"github.com/fatih/color"
)

func mercuCheckHeader(force bool, yearFlag string, authorFlag string, suffixArr []string) error {
		rootDir, err := os.Getwd()
  if err != nil {
    return err
  }

  err = filepath.WalkDir(rootDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		//check if folder or file should be skipped
		if shouldSkipDirOrFile(info.Name(), info.IsDir()) {
			if info.IsDir() {
				color.Cyan("skipped: %s tree", info.Name())
				return filepath.SkipDir
			}
			color.Cyan("skipped: %s file", info.Name())
			return nil
		}

		suffix := filepath.Ext(path)
		suffixFlag := flag.Lookup("suffix")

		if !contains(suffix, suffixArr) && suffixFlag.Value.String() != "" {
			return nil
		}

		var templateContent string
		//get correct template for this suffix
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
		// relative, err := filepath.Rel(rootDir, path)
		// if err != nil {
		// 	fmt.Println("error with relative path")
		// 	return err
		// }
		templateLinesLen := len(strings.Split(templateContent, "\n"))
    fileName := filepath.Base(path)
		filenameModded := "'" + fileName + "'"
		var trimmedYearRange string
		cmd := exec.Command("hg", "log", "--template", "{date|shortdate}\n", "-r", "reverse(ancestors(file("+filenameModded+")))")
		dir := filepath.Dir(path)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error running 'hg log years' command for file %s: %v\nOutput: %s\n", path, err, output)
			trimmedYearRange = yearFlag
		} else {
			commitDates := strings.Fields(string(output))
			var years []string
			for _, date := range commitDates {
				years = append(years, date[:4])
			}
			years = getUniques(years)
			yearRange := formatYearRange(years)
			trimmedYearRange = strings.ReplaceAll(yearRange, `"`, "")
		}

		var trimmedAuthorList string
		var authors []string
		cmd2 := exec.Command("hg", "log", "--template", "{author|person} <{author|email}>\n", path)
		cmd2.Dir = dir
		output2, err := cmd2.Output()
		if err != nil {
			fmt.Printf("Error running 'hg log authors' command for file %s: %v\nOutput: %s\n", path, err, output2)
			trimmedAuthorList = authorFlag
		} else {
			authors = getUniques(strings.Split(strings.TrimSpace(string(output2)), "\n"))
      sort.Sort(sort.Reverse(sort.StringSlice(authors)))
			authorList := strings.Join(authors, "\n *           ")
			trimmedAuthorList = strings.ReplaceAll(authorList, `"`, "")
		}
		//Modify basic template with correct information
		templateContent = strings.ReplaceAll(templateContent, "{YEARS}", trimmedYearRange)
		templateContent = strings.ReplaceAll(templateContent, "{AUTHOR}", trimmedAuthorList)
		templateContent = strings.ReplaceAll(templateContent, "{FILENAME}", filepath.Base(path))

		existingContent, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			return err
		}

		//check for existing header. in max template length + authors num. lines if 3 authors get removed this will result in error of not finding header
		maxLines := templateLinesLen + len(authors) + 3
		existingHeader := ""

		headerLinesSplit := strings.Split(string(existingContent), "\n")
		//check if file has atleast enough lines to contain template
		if len(headerLinesSplit) < maxLines {
			maxLines = len(headerLinesSplit)
		}
		switch {
		case contains(suffix, defaultSuffix):
			for i := 0; i < maxLines; i++ {
				line := headerLinesSplit[i]
				if strings.Contains(line, "************************************************************/") {
					headerStartIndex := strings.Index(string(existingContent), "****************************************************************/")
					if headerStartIndex != -1 {
						existingHeader = string(existingContent[:headerStartIndex+len("****************************************************************/")])
						existingContent = existingContent[headerStartIndex+len("****************************************************************/"):]
						break
					}
				}
			}
		case contains(suffix, pySuffix):
			for i := 0; i < maxLines; i++ {
				line := headerLinesSplit[i]
				if strings.Contains(
					line,
					`****************************************************"""`,
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
						break
					}
				}
			}
		case contains(suffix, htmlSuffix):
			for i := 0; i < maxLines; i++ {
				line := headerLinesSplit[i]
				if strings.Contains(
					line,
					`------------------------------------------------------->`,
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
						break
					}
				}
			}
		default:
			fmt.Println("error no suffix found")
			return nil
		}
		//clean useless empty space and linebreaks
		cleanedHeader := cleanString(existingHeader)
		cleanedtemplateContent := cleanString(templateContent)
		if cleanedHeader == cleanedtemplateContent {
			fmt.Printf("File %s is correct \n", path)
			return nil
		}
		if !force && len(existingHeader) < 10 {
			color.Red("No header found: %s \n", path)
			return nil
		}
		//compare lines and show difference
		oldLines := strings.Split(existingHeader, "\n")
		// newLines := strings.Split(templateContent, "\n")
		if !force {
			// if previosly found header but the header is smaller than template we assume it was not correct header
			if len(oldLines) < templateLinesLen {
				color.Red("No centria copyright header found: %s \n!\n", path)
				return nil
			}
			color.Red("file %s needs fix \n \n", path)
			// err = showDifferences(newLines, oldLines, templateLinesLen)
			if err != nil {
				color.Red("error with file %s check manually or consider ignoring if forcing header.\n!\n! ", path)
			}
			return nil
		}
		var newContent string
		//_ = showDifferences(newLines, oldLines)
		if len(oldLines) < templateLinesLen {
			//here we assume header was wrong and force new at beginning
			existingContent2, _ := os.ReadFile(path)
			newContent = templateContent + "\n" + string(existingContent2)
		} else {
			// Combine the new header with the existing content
			newContent = templateContent + "\n" + string(existingContent)
		}
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
