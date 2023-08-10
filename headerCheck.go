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

func gitCheckHeader(rootDir string, force bool, yearFlag string, authorFlag string) error {
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
		//add more suffix/templates when needed
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
		cmd := exec.Command("git", "log", "--follow", "--reverse", "--pretty=format:\"%as\"", "--", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
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
		cmd2 := exec.Command("git", "log", "--follow", "--reverse", "--pretty=format:\"%an <%ae>\"", "--", path)
		output2, err := cmd2.Output()
		if err != nil {
			trimmedAuthorList = authorFlag
		} else {
			authors := deduplicateAndSort(strings.Split(strings.TrimSpace(string(output2)), "\n"))
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

		switch {
		case contains(suffix, defaultSuffix):
			headerStartIndex := strings.Index(string(existingContent), "**********************************************************/")
			if headerStartIndex != -1 {
				existingHeader = string(existingContent[:headerStartIndex+len("**********************************************************/")])
				existingContent = existingContent[headerStartIndex+len("**********************************************************/"):]
			}
		case contains(suffix, pySuffix):
			headerStartIndex := strings.Index(string(existingContent), `********************************************************"""`)
			if headerStartIndex != -1 {
				existingHeader = string(existingContent[:headerStartIndex+len(`********************************************************"""`)])
				existingContent = existingContent[headerStartIndex+len(`********************************************************"""`):]
			}
		case contains(suffix, htmlSuffix):
			headerStartIndex := strings.Index(string(existingContent), `---------------------------------------------------------->`)
			if headerStartIndex != -1 {
				existingHeader = string(existingContent[:headerStartIndex+len(`---------------------------------------------------------->`)])
				existingContent = existingContent[headerStartIndex+len(`---------------------------------------------------------->`):]
			}
		default:
			fmt.Println("error no suffix found")
			return nil
		}

		if existingHeader == templateContent {
			return nil
		}
		if !force {
			color.Red("file %s needs fix \n  \n", path)
			oldLines := strings.Split(existingHeader, "\n")
			newLines := strings.Split(templateContent, "\n")
			//fmt.Printf("old lines len %d : new lines %d", len(oldLines), len(newLines))
			for i := 0; i < len(oldLines) && i < len(newLines); i++ {
				if oldLines[i] != newLines[i] {
					color.Red("Line: %d - %s", i+1, oldLines[i])
					color.Green("Line: %d + %s", i+1, newLines[i])
					fmt.Println()
				}
			}
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

		fmt.Printf("Copyright header fixed for file: %s\n", path)
		return nil
	})
	return err
}
