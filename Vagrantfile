# -*- mode: ruby -*-
# vi: set ft=ruby :

unless Vagrant.has_plugin?("vagrant-docker-compose")
  raise "Install plugin vagrant-docker-compose with command $> vagrant plugin install vagrant-docker-compose"
end


github_api_token = if ENV['GITHUB_API_TOKEN'].nil? || ENV['GITHUB_API_TOKEN'].empty? then
                     "7bf553ea09a665829455afd0f0541342fa85d71b"
                   else
                     ENV['GITHUB_API_TOKEN']
                   end

github_organization = if ENV['GITHUB_ORGANIZATION'].nil? || ENV['GITHUB_ORGANIZATION'].empty? then
                        "intervals-mining-lab"
                       else
                        ENV['GITHUB_ORGANIZATION']
                      end

github_team = if ENV['GITHUB_TEAM'].nil? || ENV['GITHUB_TEAM'].empty? then
                        "libiada-developers"
                      else
                        ENV['GITHUB_TEAM']
                      end


# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://atlas.hashicorp.com/search.
  config.vm.box = "ubuntu/trusty64"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  config.vm.network "private_network", ip: "192.168.33.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "./", "/vagrant_data"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  # config.vm.provider "virtualbox" do |vb|
  #   # Display the VirtualBox GUI when booting the machine
  #   vb.gui = true
  #
  #   # Customize the amount of memory on the VM:
  #   vb.memory = "1024"
  # end
  #
  # View the documentation for the provider you are using for more
  # information on available options.



  config.vm.provision :docker

  ##  Docker on vagrant have issues with dns. So just use google public dns.
  config.vm.provision "shell", inline: <<-SHELL
    sed -i -- 's/#DOCKER_OPTS="--dns 8.8.8.8 --dns 8.8.4.4"/DOCKER_OPTS="--dns 8.8.8.8 --dns 8.8.4.4"/g' /etc/default/docker
    service docker restart

    printf "#!/bin/bash \n /usr/bin/docker run -e GITHUB_API_TOKEN=#{github_api_token} -e GITHUB_ORGANIZATION=#{github_organization} -e GITHUB_TEAM=#{github_team} vagrant_github-authorized-keys authorize \\$1" > /etc/ssh_auth

    chmod +x /etc/ssh_auth

    if  ! grep -Fq  'AuthorizedKeysCommand /etc/ssh_auth' /etc/ssh/sshd_config;
    then
      printf "\nAuthorizedKeysCommand /etc/ssh_auth\n" >> /etc/ssh/sshd_config
    fi

    if  ! grep -Fq  'AuthorizedKeysCommandUser root' /etc/ssh/sshd_config;
    then
      printf "\nAuthorizedKeysCommandUser root\n" >> /etc/ssh/sshd_config
    fi

    service ssh restart
  SHELL

  config.vm.provision :docker_compose,
                      yml: "/vagrant/docker-compose-vagrant.yaml",
                      env: {
                          GITHUB_API_TOKEN: github_api_token,
                          GITHUB_ORGANIZATION: github_organization,
                          GITHUB_TEAM: github_team
                      }
end
