package configparser

import "errors"

var (
	// ErrTargetNotPointer godoc
	ErrTargetNotPointer = errors.New("target is not pointer")
	// ErrFileDoesNotExist godoc
	ErrFileDoesNotExist = errors.New("file does not exist")
)
