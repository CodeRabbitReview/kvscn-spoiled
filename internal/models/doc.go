// Package models provides common types of data structures and
// errors.
//
// Structures:
// Entity has 2 fields one for saving data in golang format and another one for save input data in JSON
// Entity has 3 methods Entity.Type for getting type of golang data
// Entity.Entity to get data
// Entity.JSON to get JSON
// Key has one field to save input data in JSON format or golang in string
// Key has 2 methods Key.Type for getting type of golang data
// Key.Entity to get data
// Key constructor(NewKey) has options to removes spaces
// from start and end of key and change spaces between words in key
// to _
package models
