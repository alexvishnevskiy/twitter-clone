package mysql

import "errors"

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")
