GO_LDFLAGS = -ldflags "-s -w"
GO_VERSION = 1.18

init:
	apt-get install -y ffmpeg
	apt-get install -y gcc libgl1-mesa-dev xorg-dev
	go mod download -x
	go mod tidy

cmd: init
	go build -o dion-frout $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/cmd
isglb: init
	go build -o isglb $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/isglb
sxu: init
	go build -o sxu $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/sxu

all: cmd isglb sxu

clean:
	rm -f dion-frout

docker-build:
	docker run --name dion-builder --rm -v $(pwd):/dion-frout golang:1.18 'apt-get update && apt-get install make && cd dion-frout && make all'