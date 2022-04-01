package main

import (
	"fmt"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"github.com/mishaprokop4ik/storage/internal/storage/models"
	"log"
)

var s = storage.NewStorage()

func main() {
	err := s.Put(storage.Pair{
		Key:    models.NewKey("number"),
		Entity: models.NewEntity(123),
	})
	if err != nil {
		log.Fatal(err)
	}

	v, err := s.Get(models.NewKey("number"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v.Entity())
}
