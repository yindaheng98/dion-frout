package controller

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/cloudwebrtc/nats-discovery/pkg/discovery"
	pb "github.com/yindaheng98/dion/proto"
	"log"
	"time"
)

type Client interface {
	Connect(nats_addr string) error
	GetNodes() map[string]discovery.Node
	SwitchNode(id string)
	SwitchSession(session *pb.ClientNeededSession)
}

func Control(w fyne.Window, cli Client) {
	addr := GetNatsAddr(w)
	log.Println("Connecting: ", addr)
	if err := cli.Connect(addr); err != nil {
		log.Fatalln(err)
	}
	ShowNodes(w, cli)
}

func GetNatsAddr(w fyne.Window) string {
	addr := "nats://127.0.0.1:4222"
	label := widget.NewLabel("NATS Address:")
	input := widget.NewEntry()
	input.SetPlaceHolder(addr)
	addrCh := make(chan string)
	connect := widget.NewButton("Connect!", func() {
		if input.Text != "" {
			addr = input.Text
		}
		addrCh <- addr
	})
	w.SetContent(container.NewVBox(label, input, connect))
	return <-addrCh
}

func ShowNodes(w fyne.Window, cli Client) {
	log.Println("Showing nodes")
	for {
		log.Println("Updating node list")
		var ids []string
		var nodes []discovery.Node
		for id, node := range cli.GetNodes() {
			ids = append(ids, id)
			nodes = append(nodes, node)
		}
		list := widget.NewList(
			func() int {
				return len(nodes)
			},
			func() fyne.CanvasObject {
				return widget.NewButton("", func() {})
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				button := o.(*widget.Button)
				button.SetText(ids[i])
				button.OnTapped = func() {
					cli.SwitchNode(ids[i])
				}
			})
		w.SetContent(list)
		<-time.After(time.Second)
	}
}
