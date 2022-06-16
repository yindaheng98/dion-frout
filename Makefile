GO_LDFLAGS = -ldflags "-s -w"
GO_VERSION = 1.18

download:
	go mod download -x
	go mod tidy
islb: download
	go build -o bin/islb $(GO_LDFLAGS) github.com/yindaheng98/dion/cmd/islb
stupid: download
	go build -o bin/stupid $(GO_LDFLAGS) github.com/yindaheng98/dion/cmd/stupid
isglb: download
	go build -o bin/isglb $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/cmd/aliyun/isglb
sxu: download
	go build -o bin/sxu $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/cmd/aliyun/sxu

server: islb stupid isglb sxu

init:
	apt-get install -y ffmpeg
	apt-get install -y gcc libgl1-mesa-dev xorg-dev

cmd: init download
	go build -o bin/dion-frout $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/cmd

all: cmd server

clean:
	rm -rf bin/

PROXY = http://192.168.1.2:10801
docker-build:
	docker run --rm -e HTTP_PROXY=$(PROXY) -e HTTPS_PROXY=$(PROXY) -v $(shell pwd):/dion-frout golang:1.18-buster sh -c 'apt update && apt install -y make && cd /dion-frout && make all'
docker-build-server:
	docker run --rm -e HTTP_PROXY=$(PROXY) -e HTTPS_PROXY=$(PROXY) -v $(shell pwd):/dion-frout golang:1.18-buster sh -c 'apt update && apt install -y make && cd /dion-frout && make server'