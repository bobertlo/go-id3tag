package id3

import (
	"errors"
)

var (
	ErrEOF     = errors.New("EOF")
	ErrNoTag   = errors.New("no tag")
	ErrInvalid = errors.New("invalid tag")
	ErrRead    = errors.New("could not read file")
)
