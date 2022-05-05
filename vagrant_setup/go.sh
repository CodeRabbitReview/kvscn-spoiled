#!/bin/bash

sudo curl -O https://storage.googleapis.com/golang/go1.18.1.linux-amd64.tar.gz
sudo tar -xvf go1.18.1.linux-amd64.tar.gz
sudo mv go /usr/local
sudo echo 'export GOROOT=/usr/local/go' >> /etc/profile
sudo echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
sudo echo 'export GOPATH=/usr/local' >> /etc/profile
sudo echo 'export PATH=$PATH:$GOPATH/bin' >> /etc/profile
sudo echo 'export GO111MODULE=on' >> /etc/profile
sudo echo 'export CGO_ENABLE=0' >> /etc/profile
source /etc/profile
sudo rm go1.18.1.linux-amd64.tar.gz
go install honnef.co/go/tools/cmd/staticcheck@latest
sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.45.0
cd /workspace/storage && go mod download