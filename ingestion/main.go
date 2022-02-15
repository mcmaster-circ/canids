// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package main

import (
	"fmt"
	"os"

	"github.com/mcmaster-circ/canids-v2/ingestion/engine"
)

func main() {
	if err := engine.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
