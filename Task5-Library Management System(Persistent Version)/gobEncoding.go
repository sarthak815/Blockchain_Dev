package main

import "bytes"

func (v Books) MarshalBinaryBook() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	return b.Bytes(), nil
}
