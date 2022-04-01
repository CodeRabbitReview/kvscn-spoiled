package models

import "errors"

var ErrNoSuchKey = errors.New("no such value by this key")
var ErrNilInput = errors.New("nil error in key data")
var ErrEmptyKeyString = errors.New("empty key value")
