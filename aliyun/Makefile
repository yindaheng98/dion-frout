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
	go build -o bin/isglb $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/aliyun/isglb
sxu: download
	go build -o bin/sxu $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/aliyun/sxu

all: islb stupid isglb sxu

clean:
	rm -rf bin/

PROXY = http://192.168.1.2:10801
docker-build:
	docker run --rm -e HTTP_PROXY=$(PROXY) -e HTTPS_PROXY=$(PROXY) -v $(shell pwd):/dion-frout golang:1.18-buster sh -c 'apt update && apt install -y make && cd /dion-frout && make all'