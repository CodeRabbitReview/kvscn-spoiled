# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
    config.vm.provider "virtualbox" do |v|
      v.memory = 2048
      v.cpus = 4
    end
    config.vm.box = "ubuntu/xenial64"

        config.vm.hostname = "mylocal.dev"
        config.vm.synced_folder ".", "/workspace/storage/",
          :owner => 'vagrant',
          :group => 'root'
        config.vm.provision :shell, :path => "vagrant_setup/init.sh"
        config.vm.provision :shell, :path => "vagrant_setup/docker.sh"
        config.vm.provision :shell, :path => "vagrant_setup/go.sh"
        config.vm.provision :shell, :path => "vagrant_setup/gradle.sh"

        config.vm.network "forwarded_port", guest: 8080, host: 8080,
            auto_correct: true
end
