package httpserver

import (
	"context"
	"crypto/tls"
	zlog "github.com/mishaprokop4ik/storage/internal/log"
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
	server   *http.Server
	certPath string
	keyPath  string
}

type KeyCertPaths struct {
	Key         string
	Certificate string
}

// NewHTTPServer is a constructor of HTTPServer
func NewHTTPServer(h http.Handler, certPaths ...KeyCertPaths) *HTTPServer {
	cfg := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		},
	}
	for _, c := range certPaths {
		pair, err := tls.LoadX509KeyPair(c.Certificate, c.Key)
		if err != nil {
			zlog.Log.Error(err, "can not read certificate", "key file name: ", c.Key,
				"cert file name:", c.Certificate)
			return nil
		}
		cfg.Certificates = append(cfg.Certificates, pair)
	}
	return &HTTPServer{server: &http.Server{
		Addr:           ":8080",
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 0,
		TLSConfig:      cfg,
		TLSNextProto:   make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}}
}

type resumer interface {
	SendRecovered(addr string)
}

// Run runs https server and take recovered data and send it
// to server in parallel.
// Run catch system signal and display it
func (s *HTTPServer) Run(r resumer) {
	go func() {
		if err := s.server.ListenAndServeTLS(s.certPath,
			s.keyPath); err != nil {
			zlog.Log.Error(err, "can not start https server")
			return
		}
	}()

	go r.SendRecovered(s.server.Addr)
	zlog.Log.Info("server is running", "on", s.server.Addr)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	signal.Notify(sc, syscall.SIGTERM)
	sig := <-sc
	zlog.Log.Info("caught system", "signal", sig)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.server.Shutdown(tc)
	cancel()
	zlog.Log.WithName("storage").Info("server stopped", "time", time.Now().String())
}
