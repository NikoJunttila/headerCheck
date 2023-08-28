/****************************************************************
 *
 *  File   : findTemplate.go
 *  Author : 
 *
 *  Copyright (C)  Centria University of Applied Sciences.
 *  All rights reserved.
 *
 *  Unauthorized copying of this file, via any medium is strictly
 *  prohibited.
 *
 ****************************************************************/
package main

import (
  "fmt"
  "os"
  "flag"
  "path/filepath"
)

func flagTemplate()(bool, string){
		templateFlag := flag.Lookup("template")
    flagVal := templateFlag.Value.String()
    fmt.Println(flagVal)
		if flagVal != "" {
      if _, err := os.Stat(flagVal); err == nil{
		    customTemplate, err := os.ReadFile(flagVal)
        if err != nil {
          return false, ""
        }
        return true, string(customTemplate)
      }
		}
		return false, ""
}

func getGwdTemplate()(bool, string){
    customTempPath := "template.txt"
		customTemplate, err := os.ReadFile(customTempPath)
	  if err != nil {
      return false, ""
		} else {
      return true, string(customTemplate)
    }
}

func getGlobalTemplate()(bool, string){
    exePath, err := os.Executable()
    if err != nil {
      return false, ""
    }
    tempPath := "template.txt"
    absolutePath := filepath.Join(filepath.Dir(exePath), tempPath)
    customTemplate, err := os.ReadFile(absolutePath)
    if err != nil {
      return false, ""
    }   
    return true, string(customTemplate)
}
