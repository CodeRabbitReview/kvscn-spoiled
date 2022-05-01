#!/bin/bash

sudo apt install apt-transport-https ca-certificates software-properties-common -y
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt update -y
sudo apt-cache policy docker-ce
sudo apt install docker-ce -y
sudo usermod -aG docker $(whoami)
sudo chmod 666 /var/run/docker.sock

sudo touch /var/lib/cloud/instance/locale-check.skip
