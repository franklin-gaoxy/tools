#!/bin/bash

## describe: Suitable for debian12 system
## 
##

function vars() {
    export remove_log = ">>/dev/null 2>&1"
}

function print_red() {
    echo -e "\033[31m$1\033[0m"
}

function print_green() {
    echo -e "\033[32m$1\033[0m"
}

function print_yellow() {
    echo -e "\033[33m$1\033[0m"
}

function print_bullet() {
    echo -e "\033[34m$1\033[0m"
}

function check() {
    if [[ $? == 0 ]];then
        print_green "Success: $1"
    else
        print_red "Error: $2"
        exit 1
    fi
}

# Install Docker code start

function uninstall_docker() {
    for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do apt-get remove $pkg; done
}

function install_repository() {
    # Add Docker's official GPG key:
    apt-get update
    check "update successfully" "Failed to update apt repositories. Please check your network connection."
    print_bullet "install tools ..."
    apt-get install -y ca-certificates curl

    # check dns hostname
    curl `hostname` $remove_log
    if [[ $? != 0 ]]; then
        print_yellow "Warning: DNS hostname resolution failed. Adding hostname to /etc/hosts."
        hostname=`hostname`
        sed -i "/${hostname}/d" /etc/hosts
        echo "127.0.0.1 ${hostname}" >> /etc/hosts
    fi

    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    echo \
"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
$(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
    tee /etc/apt/sources.list.d/docker.list > /dev/null
    apt-get update
}

function install_docker() {
    print_bullet "start install docker ..."
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-compose
    check "docker install successfully" "docker install failed"
}

function configure_docker() {
    # add config
    mkdir -p /etc/docker
    cat <<EOF >/etc/docker/daemon.json
{
  "registry-mirrors": [
    "https://registry.cn-hangzhou.aliyuncs.com",
    "https://mxov9ds4.mirror.aliyuncs.com",
    "https://mirror.ccs.tencentyun.com",
    "https://docker.sunzishaokao.com",
    "https://docker.xuanyuan.me/",
    "https://hub-mirror.c.163.com"
  ],
  "data-root": "/data/docker-data",
  "default-address-pools": [
    {"base": "10.20.0.0/16", "size": 24}
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "50m",
    "max-file": "5",
    "labels": "production"
  },
  "storage-driver": "overlay2",
  "storage-opts": [
  ],
  "dns": ["8.8.8.8", "1.1.1.1"],
  "experimental": true,
  "metrics-addr": "0.0.0.0:9323",
  "live-restore": true
}
EOF
}

function start_docker() {
    print_bullet "start docker service ..."
    systemctl enable docker
    systemctl start docker
    check "docker service started successfully" "Failed to start Docker service"
}

function main() {
    vars
    uninstall_docker
    install_repository
    install_docker
}
main