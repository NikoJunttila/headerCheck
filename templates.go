package main

type Template struct {
	Suffix string
	Header string
}

var templates = []Template{
	{".cpp", `/****************************************************************
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
	{".go", `/****************************************************************
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
	{".py", `"""****************************************************************
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
****************************************************************"""`},
	{".c", `/****************************************************************
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
	{".js", `/****************************************************************
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
	{".ts", `/****************************************************************
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
	{".cs", `/****************************************************************
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
	{".java", `/****************************************************************
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
	{".rs", `/****************************************************************
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
	{".qlm", `/****************************************************************
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
	{".css", `/****************************************************************
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
	{".html", `<!--------------------------------------------------
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
------------------------------------------------------------->`},
}
