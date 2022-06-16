GO_LDFLAGS = -ldflags "-s -w"
GO_VERSION = 1.18

init:
	apt-get install -y ffmpeg
	apt-get install -y gcc libgl1-mesa-dev xorg-dev
	go mod download -x
	go mod tidy

cmd: init
	go build -o bin/dion-frout $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/cmd
isglb: init
	go build -o bin/isglb $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/cmd/aliyun/isglb
sxu: init
	go build -o bin/sxu $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/cmd/aliyun/sxu

all: cmd bin/isglb bin/sxu

clean:
	rm -rf bin/

PROXY = http://192.168.1.2:10801
docker-build:
	docker run --rm -e HTTP_PROXY=$(PROXY) -e HTTPS_PROXY=$(PROXY) -v `pwd`:/dion-frout golang:1.18-buster sh -c 'apt update && apt install -y make && cd /dion-frout && make all'