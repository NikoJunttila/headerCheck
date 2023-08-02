package main

type Template struct {
	Suffix string
	Header string
}

var templates = []Template{
	{".cpp", `/****************************************************************
*
* File : {FILENAME}
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
	{".go", `/****************************************************************
*
* File : {FILENAME}
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
	{".py", `"""****************************************************************
*
* File : {FILENAME}
* Author : {AUTHOR}
* 
*
* Copyright (C) {YEARS} Centria University of Applied Sciences.
* All rights reserved.
*
* Unauthorized copying of this file, via any medium is strictly
* prohibited.
*
****************************************************************"""`},
	{".c", `/****************************************************************
*
* File : {FILENAME}
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

	// Add more templates for other file types if needed
}
