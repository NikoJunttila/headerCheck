package main

type Template struct {
	Suffix string
	Header string
}

var defaultSuffix = []string{".go", ".cpp", ".c", ".h", ".hpp", ".js", ".ts",
	".cs", ".java", ".rs", ".qlm", ".css", ".qss"}
var pySuffix = []string{".py"}
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
