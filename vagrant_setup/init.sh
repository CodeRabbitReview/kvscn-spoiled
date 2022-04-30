#!/bin/bash

sudo apt update
sudo apt install unzip

wget https://download.java.net/java/GA/jdk17.0.2/dfd4a8d0985749f896bed50d7138ee7f/8/GPL/openjdk-17.0.2_linux-x64_bin.tar.gz
tar xvf openjdk-17.0.2_linux-x64_bin.tar.gz
sudo mv jdk-17.0.2/ /opt/jdk-17/

sudo echo 'export JAVA_HOME=/opt/jdk-17' >> /etc/profile
sudo echo 'export PATH=$PATH:$JAVA_HOME/bin' >> /etc/profile
#echo 'export JAVA_HOME=/opt/jdk-17' | tee -a ~/.bashrc
#echo 'export PATH=$PATH:$JAVA_HOME/bin '|tee -a ~/.bashrc
wget https://services.gradle.org/distributions/gradle-7.4.2-bin.zip -P /tmp
sudo unzip -d /opt/gradle /tmp/gradle-7.4.2-bin.zip
sudo ln -s /opt/gradle/gradle-7.4.2 /opt/gradle/latest

sudo touch /etc/profile.d/gradle.sh

sudo echo 'export GRADLE_HOME=/opt/gradle/latest' >> /etc/profile.d/gradle.sh
sudo echo 'export PATH=${GRADLE_HOME}/bin:${PATH}' >> /etc/profile.d/gradle.sh

sudo echo 'export GRADLE_HOME=/opt/gradle/latest' >> /etc/profile
sudo echo 'export PATH=${GRADLE_HOME}/bin:${PATH}' >> /etc/profile
sudo chmod +x /etc/profile.d/gradle.sh

export GOPATH=$HOME/work
sudo curl -O https://storage.googleapis.com/golang/go1.18.1.linux-amd64.tar.gz
sudo tar -xvf go1.18.1.linux-amd64.tar.gz
sudo mv go /usr/local
sudo echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
export PATH=$PATH:/usr/local/go/bin
sudo rm go1.18.1.linux-amd64.tar.gz

sudo apt install apt-transport-https ca-certificates software-properties-common -y
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt update -y
sudo apt-cache policy docker-ce
sudo apt install docker-ce -y
sudo usermod -aG docker $(whoami)
sudo chmod 666 /var/run/docker.sock

sudo rm openjdk-17.0.2_linux-x64_bin.tar.gz

sudo touch /var/lib/cloud/instance/locale-check.skip