package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// HTTPServer represents connected Server, Logger and Handler
//
// Server is a http.Server struct.
// Logger is a log.Logger struct
// Handler is an interface http.Handler
// it responds to an HTTP request
type HTTPServer struct {
	server          *http.Server
	logger          *log.Logger
	fileRecoverName string
}

// NewHTTPServer is a constructor of HTTPServer
func NewHTTPServer(l *log.Logger, h http.Handler) *HTTPServer {
	return &HTTPServer{server: &http.Server{
		Addr:           ":8080",
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 0,
	}, logger: l}
}

type resumer interface {
	SendRecovered(addr string)
}

//Run ..
func (s *HTTPServer) Run(r resumer) {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Fatal(err)
		}
	}()

	go r.SendRecovered(s.server.Addr)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	signal.Notify(sc, syscall.SIGTERM)
	sig := <-sc
	s.logger.Printf("\ncaught signal %v", sig)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.server.Shutdown(tc)
	cancel()
}
