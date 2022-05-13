#!/bin/bash

wget https://download.java.net/java/GA/jdk17.0.2/dfd4a8d0985749f896bed50d7138ee7f/8/GPL/openjdk-17.0.2_linux-x64_bin.tar.gz
tar xvf openjdk-17.0.2_linux-x64_bin.tar.gz
sudo mv jdk-17.0.2/ /opt/jdk-17/

sudo echo 'export JAVA_HOME=/opt/jdk-17' >> /etc/profile
sudo echo 'export PATH=$PATH:$JAVA_HOME/bin' >> /etc/profile
wget https://services.gradle.org/distributions/gradle-7.4.2-bin.zip -P /tmp
sudo unzip -d /opt/gradle /tmp/gradle-7.4.2-bin.zip
sudo ln -s /opt/gradle/gradle-7.4.2 /opt/gradle/latest

sudo touch /etc/profile.d/gradle.sh

sudo echo 'export GRADLE_HOME=/opt/gradle/latest' >> /etc/profile.d/gradle.sh
sudo echo 'export PATH=${GRADLE_HOME}/bin:${PATH}' >> /etc/profile.d/gradle.sh

sudo echo 'export GRADLE_HOME=/opt/gradle/latest' >> /etc/profile
sudo echo 'export PATH=$GRADLE_HOME/bin:${PATH}' >> /etc/profile
sudo chmod +x /etc/profile.d/gradle.sh

sudo rm openjdk-17.0.2_linux-x64_bin.tar.gz
