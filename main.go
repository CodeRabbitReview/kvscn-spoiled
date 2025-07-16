package main

import (
	httpserver "github.com/mishaprokop4ik/storage/internal/http_server"
	"github.com/mishaprokop4ik/storage/internal/http_server/handlers"
	"github.com/mishaprokop4ik/storage/internal/recoverer"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"log"
	"os"
)

func main() {
	log := log.New(os.Stdout, "storage", log.LstdFlags)
	r := recoverer.NewTransactionLogger(recoverer.DefaultSaveFile, log)

	storage := storage.NewStorage(r)
	server := httpserver.NewHTTPServer(log, handlers.NewStorage(log, storage))
	server.Run(r)
}
