package httpserver

import (
	"log"
	"net/http"
	"time"
)

// HTTPServer represents connected Server, Logger and Handler
//
// Server is a http.Server struct.
// Logger is a log.Logger struct
// Handler is an interface http.Handler
// it responds to an HTTP request
type HTTPServer struct {
	Server  *http.Server
	Logger  *log.Logger
	Handler http.Handler
}

// NewHTTPServer is a constructor of HTTPServer
func NewHTTPServer(l *log.Logger, h http.Handler) *HTTPServer {
	return &HTTPServer{Server: &http.Server{
		Addr:           ":8080",
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 0,
	}, Logger: l, Handler: h}
}
