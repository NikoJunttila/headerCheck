package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func checkHeader(rootDir string, force bool) error {
	err := filepath.WalkDir(rootDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		// Check if the file has a recognized suffix for which we have a template
		suffix := filepath.Ext(path)
		var templateContent string
		for _, template := range templates {
			if strings.EqualFold(template.Suffix, suffix) {
				templateContent = template.Header
				break
			}
		}

		if templateContent == "" {
			return nil
		}
		// Retrieve the commit dates of the file using the "git log" command
		var trimmedYearRange string
		cmd := exec.Command("git", "log", "--follow", "--reverse", "--pretty=format:\"%as\"", "--", path)
		output, err := cmd.CombinedOutput() // Use CombinedOutput to get both stdout and stderr
		if err != nil {
			fmt.Printf("Error running 'git log' command for file %s: %v\nOutput: %s\n", path, err, output)
			trimmedYearRange = "2023"
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
			fmt.Printf("Error running 'git log' command for file %s: %v\n", path, err)
			trimmedAuthorList = "Unknown"
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

		switch suffix {
		case ".go", ".cpp", ".c", ".js", ".ts", ".cs", ".java", ".rs", ".qlm", ".css":
			headerStartIndex := strings.Index(string(existingContent), "****************************************************************/")
			if headerStartIndex != -1 {
				existingHeader = string(existingContent[:headerStartIndex+len("****************************************************************/")])
				existingContent = existingContent[headerStartIndex+len("****************************************************************/"):]
			}
		case ".py":
			headerStartIndex := strings.Index(string(existingContent), `**************************************************************"""`)
			if headerStartIndex != -1 {
				existingHeader = string(existingContent[:headerStartIndex+len(`**************************************************************"""`)])
				existingContent = existingContent[headerStartIndex+len(`**************************************************************"""`):]
			}
		case ".html":
			headerStartIndex := strings.Index(string(existingContent), `------------------------------------------------------------->`)
			if headerStartIndex != -1 {
				existingHeader = string(existingContent[:headerStartIndex+len(`------------------------------------------------------------->`)])
				existingContent = existingContent[headerStartIndex+len(`------------------------------------------------------------->`):]
			}
		default:
			fmt.Println("error no suffix found")
			return nil
		}

		if existingHeader == templateContent {
			return nil
		}
		if !force {
			fmt.Printf("file %s needs fixing \n", path)
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

func deduplicateAndSort(input []string) []string {
	uniqueMap := make(map[string]bool)
	for _, v := range input {
		uniqueMap[v] = true
	}
	var uniqueList []string
	for k := range uniqueMap {
		uniqueList = append(uniqueList, k)
	}
	sort.Strings(uniqueList)
	return uniqueList
}
func formatYearRange(years []string) string {
	if len(years) == 0 {
		return ""
	} else if len(years) == 1 {
		return years[0]
	} else {
		return years[0] + "-" + years[len(years)-1]
	}
}
