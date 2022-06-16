GO_LDFLAGS = -ldflags "-s -w"
GO_VERSION = 1.18

download:
	go mod download -x
	go mod tidy
init:
	apt-get install -y ffmpeg
	apt-get install -y gcc libgl1-mesa-dev xorg-dev

cmd: init download
	go build -o bin/dion-frout $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/cmd

all: cmd

clean:
	rm -rf bin/
