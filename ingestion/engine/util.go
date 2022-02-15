// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package engine

import "bufio"

// scanByteCounter is for counting the number of bytes read during a buffered
// file read.
type scanByteCounter struct {
	BytesRead int64 // BytesRead indicate the number of bytes read
}

// wrap is an extension of the bufio reader using a split function. It wil
// calculate the number of bytes read.
func (s *scanByteCounter) wrap(split bufio.SplitFunc) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (int, []byte, error) {
		adv, tok, err := split(data, atEOF)
		s.BytesRead += int64(adv)
		return adv, tok, err
	}
}
