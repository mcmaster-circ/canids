// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package engine

// header is Zeek TSV file header.
type header struct {
	separator    string   // separator is TSV delimeter character
	setSeperator string   // setSeperator is TSV set delimeter character
	emptyField   string   // emptyField is string identifying empty field
	unsetField   string   // unsetField is string identifying unset field
	fields       []string // fields is a list of available fields
	types        []string // types are field types
}
