package appledouble

import "errors"

var (
	ErrNotAppleDouble = errors.New("not an AppleDouble file")
	ErrAppleDoubleV1  = errors.New("pre-OSX AppleDouble files are not supported")
	ErrCorrupt        = errors.New("corrupt AppleDouble file")
)
