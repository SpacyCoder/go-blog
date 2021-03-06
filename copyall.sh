#!/bin/bash
export GOPATH=/home/spacy/go
export GOOS=linux
export CGO_ENABLED=0

cd accountservice;go get;go build -o accountservice-linux-amd64;echo built `pwd`;cd ..
cd vipservice;go get;go build -o vipservice-linux-amd64;echo built `pwd`;cd ..
cd healthchecker;go get;go build -o healthchecker-linux-amd64;echo built `pwd`;cd ..
cd imageservice;go get;go build -o imageservice-linux-amd64;echo built `pwd`;cd ..

cp healthchecker/healthchecker-linux-amd64 accountservice/
cp healthchecker/healthchecker-linux-amd64 vipservice/
cp healthchecker/healthchecker-linux-amd64 imageservice/

docker build -t spacycoder/accountservice accountservice/
docker service rm accountservice
docker service create --log-driver=gelf --log-opt gelf-address=udp://192.168.99.100:12202 --log-opt gelf-compression-type=none --name=accountservice --replicas=1 --network=my_network -p=6767:6767 spacycoder/accountservice

docker build -t spacycoder/vipservice vipservice/
docker service rm vipservice
docker service create --log-driver=gelf --log-opt gelf-address=udp://192.168.99.100:12202 --log-opt gelf-compression-type=none --name=vipservice --replicas=1 --network=my_network -p=6868:6868 spacycoder/vipservice

docker build -t spacycoder/imageservice imageservice/
docker service rm imageservice
docker service create --log-driver=gelf --log-opt gelf-address=udp://192.168.99.100:12202 --log-opt gelf-compression-type=none --name=imageservice --replicas=1 --network=my_network -p=7777:7777 spacycoder/imageservice