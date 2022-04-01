// Package storage provides options to store data in key-value database
// it has 3 action in storage: Storage.Get, Storage.Put, Storage.Delete
// Storage.Get returns an object by key and error if value by key is not found
// Storage.Put push value into storage and returns error if key or value is not valid
// Storage.Delete remove value from database
// Also this package provides some error types such as:
// ErrNoSuchKey ErrNilInput ErrEmptyKeyString
// There is Pair struct that just combines input key and value.
// API:
// err := storage.Put(Pair{
//		Key:   "simple",
//		Value: "hello there",
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
// API:
// value, err := storage.Get("simple")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
// API:
// err := storage.Delete("simple")
//
//	if err != nil {
//		log.Fatal(err)
//	}
package storage
