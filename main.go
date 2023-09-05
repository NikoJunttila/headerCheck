/****************************************************************
 *
 *  File   : main.go
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
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func main() {

	var suffixes string
  var flagTemp string
	defaultProjectPath, err := os.Getwd()
  if err != nil {
    fmt.Println(err)
    return
  }
	forceFlagPtr := flag.Bool("force", false, "actually fix files instead of just showing whats wrong")
	flag.Var((*stringSliceFlag)(&foldersToSkip), "ignore", "Specify folders/files to ignore -ignore='vendor' -ignore='node_modules'")
	flag.StringVar(&suffixes, "suffix", "", "Comma-separated list of suffixes. only goes through these files -suffix='.js,.cpp,.py'")
	flag.StringVar(&flagTemp, "template", "", "custom template location")

  newSufPtr := flag.String("newSuf", "", "Add new default suffix /* */ comment style -newSuf='.HC'")
	authorFlagPtr := flag.String("author", "default", "default author if no repo histories")
	yearFlagPtr := flag.String("year", "2023", "default year if no repo histories")
	forceVsc := flag.String("vsc", "", "force version control if no .hg file -vsc='hg'")

  helpFlag := flag.Bool("usage", false, "Show help message")
 
	flag.Parse()
  
  if *helpFlag { 
      printUsage()
      os.Exit(0)
    }

  if len(*newSufPtr) > 1{
  defaultSuffix = append(defaultSuffix, *newSufPtr)
  }

	suffixArray := strings.Split(suffixes, ",")    
	//checks for .hg file if not found errors and defaults to mercurial
	dotGitfile := filepath.Join(defaultProjectPath, ".hg")
	_, err = os.Stat(dotGitfile)
	if err == nil || *forceVsc == "hg" {
		fmt.Print("using hg")
		readIgnore(".hgignore")
		err = checkHeader(*forceFlagPtr, *yearFlagPtr, *authorFlagPtr, suffixArray,"merc")
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Print("using git")
		readIgnore(".gitignore")
		err = checkHeader(*forceFlagPtr, *yearFlagPtr, *authorFlagPtr, suffixArray,"git")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("All files checked")
	if err != nil {
		color.Red("Error scanning project: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
  exePath, err := os.Executable()
    if err != nil {
      fmt.Println("no global exe??")
    }
     fmt.Println("First checks -template flag then local directory for template.txt then \nchecks global exe location folder if all 3 return false defaults to template inside code")
     color.Red("Global Executable location: %s \n", exePath)
     fmt.Println("Usage:")
     fmt.Println("  headerCheck [options]")
     fmt.Println("\nOptions:")
     flag.VisitAll(func(f *flag.Flag) {
       fmt.Printf("  -%s: %s (default: %v)\n", f.Name, f.Usage, f.DefValue)
     })
}

// ignore folders/files stuff
type stringSliceFlag []string

func (ssf *stringSliceFlag) String() string {
	return strings.Join(*ssf, ", ")
}

func (ssf *stringSliceFlag) Set(value string) error {
	*ssf = append(*ssf, value)
	return nil
}
