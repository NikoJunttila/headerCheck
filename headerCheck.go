/****************************************************************
 *
 *  File   : headerCheck.go
 *  Author : NikoJunttila <89527972+NikoJunttila@users.noreply.github.com>
 *           Niko Junttila <niko.junttila2@centria.fi>
 *
 *  Copyright (C) 2023-2024 Centria University of Applied Sciences.
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

	"github.com/fatih/color"
)

var FixesCheck bool = false

func checkHeader(force bool, yearFlag string, authorFlag string, suffixArr []string, gitOrMerc string, verbose bool, noHeadIgn bool) error {
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}

	var templateContentBody string
	if check, templateCustom := flagTemplate(); check {
		fmt.Println("Using given template")
		templateContentBody = templateCustom
	} else if check, templateCustom = getGwdTemplate(); check {
		fmt.Println("Using template in directory")
		templateContentBody = templateCustom
	} else if check, templateCustom = getGlobalTemplate(); check {
		fmt.Println("Using global template")
		templateContentBody = templateCustom
	} else {
		fmt.Println("Using default hardcoded template")
		templateContentBody = template
	}
	templateContentBody = strings.TrimRight(templateContentBody, "\n")

	fmt.Println("checking files...")
	err = filepath.WalkDir(rootDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(rootDir, path)
		//check if folder or file should be skipped
		//if you don't want to use relational path for file skiping change relPath to info.Name()
		if shouldSkipDirOrFile(relPath, info.IsDir()) {
			if info.IsDir() {
				color.Cyan("skipped tree: %s", relPath)
				return filepath.SkipDir
			}
			color.Cyan("skipped file: %s", relPath)
			return nil
		}

		suffix := filepath.Ext(path)
		suffixFlag := flag.Lookup("suffix")
		singleFlag := flag.Lookup("single")
		singleVal := singleFlag.Value.String()

		if singleVal != "" && filepath.Base(path) != singleVal {
			return nil
		}
		if !contains(suffix, suffixArr) && suffixFlag.Value.String() != "" {
			return nil
		}

		var templateContent string
		//get correct template for this suffix
		switch {
		case contains(suffix, defaultSuffix):
			templateContent = "/****************************************************************\n" + templateContentBody + "****************************************************************/"
		case contains(suffix, pySuffix):
			templateContent = `"""*************************************************************` + "\n" + templateContentBody + `**********************************************************"""`
		case contains(suffix, htmlSuffix):
			templateContent = "<!--------------------------------------------------------------\n" + templateContentBody + "--------------------------------------------------------------->"
		case suffix == ".rb":
			templateContent = "=begin *********************************************************" + "\n" + templateContentBody + "******************************************************* =end"
		case suffix == ".lua":
			templateContent = "--[[************************************************************" + "\n" + templateContentBody + "**********************************************************]]"
		case suffix == ".ml" || suffix == ".mli":
			//OCaml
			templateContent = "(************************************************************" + "\n" + templateContentBody + "**********************************************************)"
		default:
			return nil
		}

		templateLinesLen := len(strings.Split(templateContent, "\n"))
		// Retrieve the commit dates of the file using the "git log" command
		var trimmedYearRange string
		var cmd *exec.Cmd
		if gitOrMerc == "git" {
			cmd = exec.Command(
				"git",
				"log",
				"--follow",
				// "--reverse",
				"--pretty=format:\"%as\"",
				"--",
				path,
			)
		} else {
			fileName := filepath.Base(path)
			filenameModded := "'" + fileName + "'"
			cmd = exec.Command("hg", "log", "--template", "{date|shortdate}\n", "-r", "reverse(ancestors(file("+filenameModded+")))")
		}
		dir := filepath.Dir(path)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			trimmedYearRange = yearFlag
		} else if string(output) == "" {
			trimmedYearRange = yearFlag
		} else {
			dates := strings.Split(string(output), "\n")
			// Reverse the order of the dates
			for i, j := 0, len(dates)-1; i < j; i, j = i+1, j-1 {
				dates[i], dates[j] = dates[j], dates[i]
			}
			//commitDates := strings.Fields(string(output))
			var years []string
			for _, date := range dates {
				if gitOrMerc == "git" {
					years = append(years, date[:5])
				} else {
					years = append(years, date[:4])
				}
			}
			years = getUniques(years)
			yearRange := formatYearRange(years)
			trimmedYearRange = strings.ReplaceAll(yearRange, `"`, "")
		}

		var trimmedAuthorList string
		var authors []string
		// Retrieve the authors of the file using the "git log" command
		var cmd2 *exec.Cmd
		if gitOrMerc == "git" {
			cmd2 = exec.Command(
				"git",
				"log",
				"--follow",
				// "--reverse",
				"--pretty=format:\"%an <%ae>\"",
				"--",
				path,
			)
		} else {
			cmd2 = exec.Command("hg", "log", "--template", "{author|person} <{author|email}>\n", path)
		}
		cmd2.Dir = dir
		output2, err := cmd2.Output()
		if err != nil {
			if verbose {
				fmt.Printf("Error running 'git/hg log authors' command for file %s: %v\nOutput: %s\n", path, err, output2)
			}
			trimmedAuthorList = authorFlag
		} else if string(output2) == "" {
			if verbose {
				fmt.Println("error using git. Using default author: ", authorFlag)
			}
			trimmedAuthorList = authorFlag
		} else {
			authors := strings.Split(strings.TrimSpace(string(output2)), "\n")
			// Reverse the order of the authors
			for i, j := 0, len(authors)-1; i < j; i, j = i+1, j-1 {
				authors[i], authors[j] = authors[j], authors[i]
			}
			authors = getUniques(authors)
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
				if strings.Contains(line, "*********************************************************/") {
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
		case suffix == ".rb":
			for i := 0; i < maxLines; i++ {
				line := headerLinesSplit[i]
				if strings.Contains(
					line,
					`************************************ =end`,
				) {
					headerStartIndex := strings.Index(
						string(existingContent),
						`************************************ =end`,
					)
					if headerStartIndex != -1 {
						existingHeader = string(
							existingContent[:headerStartIndex+len(`************************************ =end`)],
						)
						existingContent = existingContent[headerStartIndex+len(`************************************ =end`):]
						break
					}
				}
			}
		case suffix == ".lua":
			for i := 0; i < maxLines; i++ {
				line := headerLinesSplit[i]
				if strings.Contains(
					line,
					`************************************************]]`,
				) {
					headerStartIndex := strings.Index(
						string(existingContent),
						`************************************************]]`,
					)
					if headerStartIndex != -1 {
						existingHeader = string(
							existingContent[:headerStartIndex+len(`************************************************]]`)],
						)
						existingContent = existingContent[headerStartIndex+len(`************************************************]]`):]
						break
					}
				}
			}
		case suffix == ".ml" || suffix == ".mli":
			for i := 0; i < maxLines; i++ {
				line := headerLinesSplit[i]
				if strings.Contains(
					line,
					`************************************************)`,
				) {
					headerStartIndex := strings.Index(
						string(existingContent),
						`************************************************)`,
					)
					if headerStartIndex != -1 {
						existingHeader = string(
							existingContent[:headerStartIndex+len(`************************************************)`)],
						)
						existingContent = existingContent[headerStartIndex+len(`************************************************)`):]
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

		oldLines := strings.Split(existingHeader, "\n")
		if noHeadIgn && len(existingHeader) < 10 {
			return nil
		}
		if noHeadIgn && len(oldLines)+1 < templateLinesLen {
			return nil
		}
		if cleanedHeader == cleanedtemplateContent {
			if verbose {
				fmt.Printf("File %s is correct \n", path)
			}
			return nil
		}

		if !force && len(existingHeader) < 10 {
			FixesCheck = true
			color.Red("No header found: %s \n", path)
			return nil
		}
		//compare lines and show difference
		newLines := strings.Split(templateContent, "\n")
		if !force {
			FixesCheck = true
			// if previosly found header but the header is smaller than template we assume it was not correct header
			if len(oldLines)+1 < templateLinesLen {
				color.Red("No centria copyright header found: %s \n!\n! \n", path)
				return nil
			}
			color.Red("\nfile %s needs fix \n", path)
			showBlockDifferences(newLines, oldLines)
			// loops through each line equally and shows differences. breaks if the template is modified between checks.
			// err = showDifferences(newLines, oldLines, templateLinesLen, authIndex)
			// if err != nil {
			// 	color.Red("error with file %s check manually or consider ignoring if forcing header.\n!\n! ", path)
			// }
			return nil
		}
		var newContent string
		if len(oldLines)+1 < templateLinesLen {
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
