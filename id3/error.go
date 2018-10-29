package id3

import (
	"errors"
)

var (
	ErrEOF           = errors.New("EOF")
	ErrInvalidHeader = errors.New("Invalid ID3 header");
	ErrRead          = errors.New("Read Error")
)
