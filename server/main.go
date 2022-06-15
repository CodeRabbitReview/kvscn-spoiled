package main

import (
	"fmt"
	httpserver "github.com/mishaprokop4ik/storage/internal/http_server"
	"github.com/mishaprokop4ik/storage/internal/http_server/handlers"
	zlog "github.com/mishaprokop4ik/storage/internal/log"
	"github.com/mishaprokop4ik/storage/internal/recoverer"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"os"
	"time"
)

func main() {
	zlog.Init("./../persistence/storage.json", "stderr")
	localCertPath := os.Getenv("LOCAL_SSL_KEY")
	if len(localCertPath) == 0 {
		zlog.Log.Error(fmt.Errorf("did not find local ssl certificate variable"), "env variable is not set")
		os.Exit(1)
	}
	localKeyPath := os.Getenv("LOCAL_SSL_CERT")
	if len(localKeyPath) == 0 {
		zlog.Log.Error(fmt.Errorf("did not find local private ssl key variable"), "env variable is not set")
		os.Exit(1)
	}
	k8sKeyPath := os.Getenv("K8S_SSL_KEY")
	if len(k8sKeyPath) == 0 {
		zlog.Log.Error(fmt.Errorf("did not find k8s ssl certificate variable"), "env variable is not set")
		os.Exit(1)
	}
	k8sCertPath := os.Getenv("K8S_SSL_CERT")
	if len(k8sKeyPath) == 0 {
		zlog.Log.Error(fmt.Errorf("did not find k8s private ssl key variable"), "env variable is not set")
		os.Exit(1)
	}
	serverULR := os.Getenv("SERVER_URL")
	if len(serverULR) == 0 {
		zlog.Log.Error(fmt.Errorf("did not find k8s private ssl key variable"), "env variable is not set")
		os.Exit(1)
	}
	zlog.Log.WithName("storage").Info("started", "time", time.Now().String())
	r := recoverer.NewTransactionLogger(recoverer.DefaultSaveFile)
	s := storage.NewStorage(r)
	server := httpserver.NewHTTPServer(handlers.NewStorage(s), 8080, httpserver.KeyCertPaths{
		Key:         localKeyPath,
		Certificate: localCertPath,
	},
		httpserver.KeyCertPaths{
			Key:         k8sKeyPath,
			Certificate: k8sCertPath,
		})
	server.URL = serverULR
	server.Run(r)
}
