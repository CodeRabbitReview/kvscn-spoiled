// Package recoverer provides options to recover data into file if server
// will unexpectedly stop.
// In next server start all data will send into server
// API:
// r := recoverer.NewTransactionLogger(recoverer.DefaultSaveFile, log)
//
//	storage := storage.NewStorage(r)
//	server := httpserver.NewHTTPServer(log, handlers.NewStorage(log, storage))
//	server.Run(r)
package recoverer
