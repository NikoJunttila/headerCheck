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
  "os/exec"
	"path/filepath"
	"strings"
  "time"
	"github.com/fatih/color"
)

func main() {
	defaultProjectPath, err := os.Getwd()
  if err != nil {
    fmt.Println(err)
    return 
  }
  usingGit := "git"
  dotHgfile := filepath.Join(defaultProjectPath, ".hg")
	_, err = os.Stat(dotHgfile)
	if err == nil{
    usingGit = "hg"
  }
  var defaultname []byte
  var defaultemail []byte
  if usingGit == "git"{
  cmd1 := exec.Command("git", "config", "--get", "user.name")
  cmd2 := exec.Command("git", "config", "--get", "user.email")
  defaultname, _ = cmd1.Output()
  defaultemail, _ = cmd2.Output()
  }else{
    cmd1 := exec.Command("hg", "config", "ui.username")
    defaultname, _ = cmd1.Output()
    cmd2 := exec.Command("hg", "config", "ui.email")
    defaultemail, _ = cmd2.Output()
  }
  
  defaultname = []byte(strings.ReplaceAll(string(defaultname), "\n", ""))
  defaultemail = []byte(strings.ReplaceAll(string(defaultemail), "\n", ""))

  defaultAuthor := fmt.Sprintf("%s <%s>", defaultname, defaultemail)
  currentYear := fmt.Sprint(time.Now().Year())
	var suffixes string
  var flagTemp string
  var single string

	forceFlagPtr := flag.Bool("force", false, "actually fix files instead of just showing whats wrong")
	flag.Var((*stringSliceFlag)(&foldersToSkip), "ignore", "Specify folders/files to ignore -ignore='vendor' -ignore='node_modules'")
	flag.StringVar(&suffixes, "suffix", "", "Comma-separated list of suffixes. only goes through these files -suffix='.js,.cpp,.py'")
	flag.StringVar(&flagTemp, "template", "", "custom template location")
  
  flag.StringVar(&single, "single", "", "If you want to only check a single file -single='my_awesome_source_file.go'")
  newSufPtr := flag.String("newSuf", "", "Add new suffix that has this comment style /* */ if not already included -newSuf='.HC'")
	authorFlagPtr := flag.String("author", defaultAuthor, "default author if no repo histories")
	yearFlagPtr := flag.String("year", currentYear, "default year if no repo histories")
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
	if usingGit == "hg" || *forceVsc == "hg" {
		fmt.Print("using hg")
		readIgnore(".hgignore")
		err = checkHeader(*forceFlagPtr, *yearFlagPtr, *authorFlagPtr, suffixArray,"merc")
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	} else {
		fmt.Print("using git")
		readIgnore(".gitignore")
		err = checkHeader(*forceFlagPtr, *yearFlagPtr, *authorFlagPtr, suffixArray,"git")
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
    
	}
  if FixesCheck{
    color.Red("Some fixes needed \n")
    os.Exit(1)
  } else{
    color.Green("All files checked and correct \n")
    os.Exit(0)
  }
}

func printUsage() {
  exePath, err := os.Executable()
    if err != nil {
      fmt.Println("no global exe??")
    }

     fmt.Println("First checks -template flag \n then directory program was executed from for template.txt then \nchecks global exe location folder if neither all 3 return false uses default template inside code")
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
