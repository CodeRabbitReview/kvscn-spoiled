package main

import (
	"fmt"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"log"
)

var s = storage.NewStorage()

func main() {
	err := s.Put(storage.Pair{
		Key:   "simple",
		Value: "hello there",
	})

	if err != nil {
		log.Fatal(err)
	}

	v, err := s.Get("simple")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)
}
