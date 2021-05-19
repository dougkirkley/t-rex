// Copyright (c) 2020-present Douglass Kirkley. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"os"
	"log"

	"github.com/dougkirkley/trex/cmd"
)

func main() {
	if err := cmd.Run(os.Args[1:]); err != nil {
	    log.Fatal(err.Error())	
	}
}
