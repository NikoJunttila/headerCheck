package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func checkHeader(path string, info os.FileInfo, err error) error {
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
	cmd := exec.Command("git", "log", "--follow", "--reverse", "--pretty=format:\"%as\"", "--", path)
	output, err := cmd.CombinedOutput() // Use CombinedOutput to get both stdout and stderr
	if err != nil {
		fmt.Printf("Error running 'git log' command for file %s: %v\nOutput: %s\n", path, err, output)
		return err
	}
	cmd2 := exec.Command("git", "log", "--follow", "--reverse", "--pretty=format:\"%an <%ae>\"", "--", path)
	output2, err := cmd2.Output()
	if err != nil {
		fmt.Printf("Error running 'git log' command for file %s: %v\n", path, err)
		return err
	}

	authors := deduplicateAndSort(strings.Split(strings.TrimSpace(string(output2)), "\n"))
	authorList := strings.Join(authors, "\n*          ")
	trimmedAuthorList := strings.ReplaceAll(authorList, `"`, "")
	commitDates := strings.Fields(string(output))

	var years []string
	for _, date := range commitDates {
		years = append(years, date[:5])
	}
	years = deduplicateAndSort(years)
	yearRange := formatYearRange(years)
	trimmedYearRange := strings.ReplaceAll(yearRange, `"`, "")

	templateContent = strings.ReplaceAll(templateContent, "{YEARS}", trimmedYearRange)
	templateContent = strings.ReplaceAll(templateContent, "{AUTHOR}", trimmedAuthorList)
	templateContent = strings.ReplaceAll(templateContent, "{FILENAME}", filepath.Base(path))

	existingContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return err
	}
	existingHeader := ""
	if suffix == ".go" || suffix == ".cpp" || suffix == ".c" || suffix == ".js" || suffix == ".ts" {
		headerStartIndex := strings.Index(string(existingContent), "****************************************************************/")
		if headerStartIndex != -1 {
			existingHeader = string(existingContent[:headerStartIndex+len("****************************************************************/")])
			existingContent = existingContent[headerStartIndex+len("****************************************************************/"):]
		}
	} else if suffix == ".py" {
		headerStartIndex := strings.Index(string(existingContent), `**************************************************************"""`)
		if headerStartIndex != -1 {
			existingHeader = string(existingContent[:headerStartIndex+len(`**************************************************************"""`)])
			existingContent = existingContent[headerStartIndex+len(`**************************************************************"""`):]
		}
	}
	if existingHeader == templateContent {
		fmt.Printf("Copyright header already exists and matches for file: %s\n", path)
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
