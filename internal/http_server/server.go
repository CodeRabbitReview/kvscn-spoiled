package httpserver

import (
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	Server   *http.Server
	Logger   *log.Logger
	Handlers http.Handler
}

func NewHTTPServer(l *log.Logger, h http.Handler) *HTTPServer {
	return &HTTPServer{Server: &http.Server{
		Addr:           ":8080",
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 0,
	}, Logger: l, Handlers: h}
}
