// Package storage provides options to store data in key-value database
// it has 4 methods in Storage: Storage.Get, Storage.Put, Storage.Delete, Storage.GetAll
// Storage.Get returns an object by key. If Storage is empty or no value by such key will be returned error
// Storage.GetAll takes all value from Storage. If no value in storage returns error
// Storage.Put push value into storage and returns error if key or entity is not valid.
// Validation checks key and entity not nil value
// Storage.Delete remove value from database if no data into storage or
// no data by this key - sends error
// There is Pair struct that just combines input key and value.
// API:
// err := Storage.Put(Pair{
//		Key:   models.NewKey("simple"),
//		Entity: models.NewEntity("hello there", nil),
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
// API:
// value, err := Storage.Get(models.NewKey("simple"))
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
// API:
// err := Storage.Delete(models.NewKey("simple"))
//
//	if err != nil {
//		log.Fatal(err)
//	}
// API:
// data, err := Storage.GetAll()
//
//	if err != nil {
//		log.Fatal(err)
//	}
package storage
