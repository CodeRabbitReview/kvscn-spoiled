package models

import "errors"

var ErrNoSuchKey = errors.New("no such value by this key")
var ErrNilInput = errors.New("nil in input data")
var ErrEmptyKey = errors.New("empty key value")
