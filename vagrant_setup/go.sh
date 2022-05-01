#!/bin/bash

export GOPATH=$HOME/work
sudo curl -O https://storage.googleapis.com/golang/go1.18.1.linux-amd64.tar.gz
sudo tar -xvf go1.18.1.linux-amd64.tar.gz
sudo mv go /usr/local
sudo echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
export PATH=$PATH:/usr/local/go/bin
sudo rm go1.18.1.linux-amd64.tar.gz
