// Package httpserver implements making http.Server with log.Logger to log data
// http.Handler in the Server that will handle http requests.
//
// To make new instance of HTTPServer you need go next steps:
// 1. Make log: log := log.New(os.Stdout, "storage", log.LstdFlags)
// 2. Make Server: server := httpserver.NewHTTPServer(log, handlers.NewStorage(log, storage))
package httpserver
