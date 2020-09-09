#!/bin/bash

#cd www
#npm run build:prod
#cd ..
#rm -rf html
#cp -R www/dist html
#
#
#go-bindata html/...

GOOS=windows GOARCH=amd64 go build -o yuque.exe *.go


#GOOS=windows GOARCH=386 go build -o enhancemposx86.exe *.go
#go build -o enhancempose *.go
#rm db2rest
#GOOS=linux GOARCH=amd64 go build -o db2rest .

#dt=$(date +%Y%m%d%H%M%S)
#docker build -t registry.cn-hangzhou.aliyuncs.com/r1/db2rest:$dt .
#docker login --username=jandyu@gmail.com -p jiangyu123 registry.cn-hangzhou.aliyuncs.com
#
#docker push registry.cn-hangzhou.aliyuncs.com/r1/db2rest:$dt
#
#sshpass -p "stock123*()" ssh root@3y.zjy8.cn "docker login --username=jandyu@gmail.com -p jiangyu123 registry.cn-hangzhou.aliyuncs.com;docker service update --image  registry.cn-hangzhou.aliyuncs.com/r1/db2rest:$dt --with-registry-auth db2rest "
#cp wssrv.exe smb://192.168.0.77/share/deploy/websocket/wssrv.exe