package main

import (
	"github.com/mishaprokop4ik/storage/internal/http_server"
	"github.com/mishaprokop4ik/storage/internal/http_server/handlers"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"log"
	"os"
)

func main() {
	log := log.New(os.Stdout, "storage", log.LstdFlags)
	storage := storage.NewStorage()
	server := httpserver.NewHTTPServer(log, handlers.NewStorage(log, storage))
	if err := server.Server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
