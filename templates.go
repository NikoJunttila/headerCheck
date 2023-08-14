/****************************************************************
*
* File   : templates.go
* Author : NikoJunttila <89527972+NikoJunttila@users.noreply.github.com>
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

type Template struct {
	Suffix string
	Header string
}

// add more languages that have /* */ comment out style
var defaultSuffix = []string{".go", ".cpp", ".c", ".h", ".hpp", ".js", ".ts",
	".cs", ".java", ".rs", ".qlm", ".css", ".qss"}

// add more languages that have """ """ python style comment out style
var pySuffix = []string{".py"}

// add more languages that have <-- --> html style comment out style
var htmlSuffix = []string{".html"}

var templates = []Template{
	{"default", `/****************************************************************
*
* File   : {FILENAME}
* Author : {AUTHOR}
*
*
* Copyright (C) {YEARS} Centria University of Applied Sciences.
* All rights reserved.
*
* Unauthorized copying of this file, via any medium is strictly
* prohibited.
*
****************************************************************/`},
	{".py", `"""*************************************************************
*
* File   : {FILENAME}
* Author : {AUTHOR}
*
*
* Copyright (C) {YEARS} Centria University of Applied Sciences.
* All rights reserved.
*
* Unauthorized copying of this file, via any medium is strictly
* prohibited.
*
*************************************************************"""`},
	{".html", `<!--------------------------------------------------------------
*
* File   : {FILENAME}
* Author : {AUTHOR}
*
*
* Copyright (C) {YEARS} Centria University of Applied Sciences.
* All rights reserved.
*
* Unauthorized copying of this file, via any medium is strictly
* prohibited.
*
--------------------------------------------------------------->`},
}
