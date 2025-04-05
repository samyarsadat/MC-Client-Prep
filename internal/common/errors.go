package common

import "errors"

// ErrDirNotExist represents an error when the destination directory does not exist.
var ErrDirNotExist = errors.New("destination directory does not exist")

// ErrFileExistsNoOverwrite represents an error when a file already exists and overwriting is disabled.
var ErrFileExistsNoOverwrite = errors.New("file already exists and overwriting is disabled")
