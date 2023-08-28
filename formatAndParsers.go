/****************************************************************
 *
 *  File   : formatAndParsers.go
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
	"sort"
	"strings"
)

func getUniques(slice []string) []string{
  var uniques []string
  for _,name := range slice{
    if !contains(name, uniques){
    uniques = append(uniques, name)
    }
  }
  return uniques
}

func formatYearRange(years []string) string {
  sort.Strings(years)
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

