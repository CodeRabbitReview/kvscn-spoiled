package main

import (
	"github.com/mishaprokop4ik/storage/internal/http_server"
	"github.com/mishaprokop4ik/storage/internal/http_server/handlers"
	storage2 "github.com/mishaprokop4ik/storage/internal/storage"
	log2 "log"
	"os"
)

func main() {
	log := log2.New(os.Stdout, "storage", log2.LstdFlags)
	storage := storage2.NewStorage()
	server := httpserver.NewHTTPServer(log, handlers.NewStorage(log, storage))
	if err := server.Server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
