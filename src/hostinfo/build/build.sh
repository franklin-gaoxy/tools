#!/bin/bash

cd ../src
go build -o hostinfo main.go


docker build -t registry.cn-hangzhou.aliyuncs.com/kubernetes_install/hostinfo:v0.0.1 .
docker push registry.cn-hangzhou.aliyuncs.com/kubernetes_install/hostinfo:v0.0.1
