package main

import (
	"sort"
)

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
