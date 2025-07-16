package main

import (
	httpserver "github.com/mishaprokop4ik/storage/internal/http_server"
	"github.com/mishaprokop4ik/storage/internal/http_server/handlers"
	zlog "github.com/mishaprokop4ik/storage/internal/log"
	"github.com/mishaprokop4ik/storage/internal/recoverer"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"time"
)

func main() {
	zlog.Log.WithName("storage").Info("started", "time", time.Now())
	r := recoverer.NewTransactionLogger(recoverer.DefaultSaveFile)
	storage := storage.NewStorage(r)
	server := httpserver.NewHTTPServer(handlers.NewStorage(storage), "localhost.pem",
		"localhost-key.pem")
	server.Run(r)
}
