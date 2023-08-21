/****************************************************************
*
* File   : formatAndParsers.go
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
	"sort"
  "strings"
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
func contains(v string, a []string) bool {
    for _, i := range a {
        if i == v {
            return true
        }
    }
    return false
}
func cleanString(needsCleaning string) string{
    needsCleaning = strings.ReplaceAll(strings.ReplaceAll(needsCleaning, "\r", ""), "\n", "")
		needsCleaning = strings.TrimSpace(needsCleaning)
  return needsCleaning
}

