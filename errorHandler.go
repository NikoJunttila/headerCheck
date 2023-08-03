/****************************************************************
*
* File   : errorHandler.go
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

import "log"

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
