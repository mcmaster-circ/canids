// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package uuid provides a unique identifier generator.
package uuid

import (
	u "github.com/satori/go.uuid"
)

// Generate will generate a new UUID.
func Generate() string {
	gen := u.NewV4()
	return gen.String()
}
