package controller

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/cloudwebrtc/nats-discovery/pkg/discovery"
	pb "github.com/yindaheng98/dion/proto"
	"log"
)

type Client interface {
	Connect(nats_addr string) error
	GetNodes() map[string]discovery.Node
	SwitchNode(id string)
	SwitchSession(session *pb.ClientNeededSession)
}

type Controller struct {
	cli Client
}

func NewController(cli Client) Controller {
	return Controller{
		cli: cli,
	}
}

func Control(a fyne.App, cli Client) {
}

func GetNatsAddr(a fyne.App) string {
	addr := "nats://127.0.0.1:4222"
	w := a.NewWindow("dion System - Please give a NATS Address")
	label := widget.NewLabel("NATS Address:")
	input := widget.NewEntry()
	input.SetPlaceHolder(addr)
	connect := widget.NewButton("Connect!", func() {
		if input.Text != "" {
			addr = input.Text
		}
		log.Println("Connecting: ", addr)
		w.Close()
	})
	w.SetContent(container.NewVBox(label, input, connect))
	w.ShowAndRun()
	return addr
}
