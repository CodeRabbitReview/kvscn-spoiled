package httpserver

import (
	"bufio"
	"context"
	"fmt"
	"github.com/mishaprokop4ik/storage/internal/client"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type recovered struct {
	method string
	data   string
}

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
func NewHTTPServer(l *log.Logger, h http.Handler, fileName string) *HTTPServer {
	return &HTTPServer{server: &http.Server{
		Addr:           ":8080",
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 0,
	}, logger: l, fileRecoverName: fileName}
}

//Run ..
func (s *HTTPServer) Run() {
	recoverData := make(chan recovered)
	c := client.NewAPI(fmt.Sprintf("http://localhost%s", s.server.Addr))
	go s.recover(recoverData)
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Fatal(err)
		}
	}()
	for {
		r := <-recoverData
		if (recovered{}) == r {
			break
		}
		go func(r recovered) {
			switch r.method {
			case http.MethodPut:
				_, err := c.AddOrUpdate(r.data)
				if err != nil {
					s.logger.Fatal(err)
				}
			case http.MethodDelete:
				_, err := c.Delete(r.data)
				if err != nil {
					s.logger.Fatal(err)
				}
			}
		}(r)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	signal.Notify(sc, syscall.SIGTERM)
	sig := <-sc
	s.logger.Printf("\ncaught signal %v", sig)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.server.Shutdown(tc)
	cancel()
}

func (s *HTTPServer) recover(c chan recovered) {
	_, err := os.Stat(s.fileRecoverName)
	if err != nil {
		s.logger.Fatal(err)
	}

	f, err := os.Open(s.fileRecoverName)
	if err != nil {
		s.logger.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		err = sc.Err()
		if err != nil {
			s.logger.Fatalf("scan file error: %v", err)
		}
		recoverData := strings.Split(sc.Text(), "\t")
		switch recoverData[0] {
		case "put":
			c <- recovered{
				method: http.MethodPut,
				data:   recoverData[1],
			}
		case "delete":
			c <- recovered{
				method: http.MethodDelete,
				data:   recoverData[1],
			}
		}
	}
	err = os.Truncate(s.fileRecoverName, 0)
	if err != nil {
		s.logger.Fatal(err)
	}
	c <- recovered{}
}
