/****************************************************************
*
* File   : headerCheck.go
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
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
  "flag"

	"github.com/fatih/color"
)

func gitCheckHeader(rootDir string, force bool, yearFlag string, authorFlag string, suffixArr []string) error {
	err := filepath.WalkDir(rootDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
//check if folder or file should be skipped
		if shouldSkipDirOrFile(info.Name(), info.IsDir()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		suffix := filepath.Ext(path)
    suffixFlag := flag.Lookup("suffix")

    if !contains(suffix, suffixArr) && suffixFlag.Value.String() != ""{
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
		// Retrieve the commit dates of the file using the "git log" command
		var trimmedYearRange string
		cmd := exec.Command(
			"git",
			"log",
			"--follow",
		  "--reverse",
			"--pretty=format:\"%as\"",
			"--",
			path,
		)
    dir := filepath.Dir(path) 
    cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
		  // fmt.Printf("Error running 'git log years' command for file %s: %v\nOutput: %s\n", path, err, output)
			trimmedYearRange = yearFlag
		} else {
			commitDates := strings.Fields(string(output))
			var years []string
			for _, date := range commitDates {
				years = append(years, date[:5])
			}
			years = deduplicateAndSort(years)
			yearRange := formatYearRange(years)
			trimmedYearRange = strings.ReplaceAll(yearRange, `"`, "")
		}

		var trimmedAuthorList string
		var authors []string
		// Retrieve the authors of the file using the "git log" command
		cmd2 := exec.Command(
			"git",
			"log",
			"--follow",
		  "--reverse",
			"--pretty=format:\"%an <%ae>\"",
			"--",
			path,
		)
    cmd2.Dir = dir
		output2, err := cmd2.Output()
		if err != nil {
		  fmt.Printf("Error running 'git log authors' command for file %s: %v\nOutput: %s\n", path, err, output2)
			trimmedAuthorList = authorFlag
		} else {
			authors = deduplicateAndSort(strings.Split(strings.TrimSpace(string(output2)), "\n"))
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
    //check for existing header. in max 14 + authors num. lines if 3 authors get removed this will result in error of not finding header
		existingHeader := ""
		maxLines := 14 + len(authors)
		headerLinesSplit := strings.Split(string(existingContent), "\n")
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
		newLines := strings.Split(templateContent, "\n")
		if !force {
    if len(oldLines) < 11{
			color.Red("No header found: %s \n", path)
      return nil
    }
			color.Red("file %s needs fix \n  \n", path)
	    err = showDifferences(newLines, oldLines)
      if err != nil {
        color.Red("error with file %s check manually or consider ignoring if forcing header.", path)
        color.Red("!\n")
        color.Red("!\n")
      }
			return nil
		}
    var newContent string
    //_ = showDifferences(newLines, oldLines)
    if len(oldLines) < 11{
      //here we assume header was wrong and force new at beginning
		  existingContent2, _ := os.ReadFile(path)
		  newContent = templateContent + "\n" + string(existingContent2) 
    } else{
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
