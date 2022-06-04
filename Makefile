GO_LDFLAGS = -ldflags "-s -w"
GO_VERSION = 1.16

init:
	apt-get install -y ffmpeg
	apt-get install -y gcc libgl1-mesa-dev xorg-dev
	go mod download -x
	go mod tidy

cmd: init
	go build -o dion-frout $(GO_LDFLAGS) github.com/yindaheng98/dion-frout/cmd

all: cmd

clean:
	rm -f dion-frout
