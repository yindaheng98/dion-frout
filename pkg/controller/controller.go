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
	ShowNodes(w, cli, SessionEntry(cli))
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

func SessionEntry(cli Client) *fyne.Container {
	user := "stupid"
	userlabel := widget.NewLabel("User:")
	userinput := widget.NewEntry()
	userinput.SetPlaceHolder(user)
	usercontainer := container.NewHBox(userlabel, userinput)
	session := "stupid"
	sessionlabel := widget.NewLabel("Session:")
	sessioninput := widget.NewEntry()
	sessioninput.SetPlaceHolder(session)
	sessioncontainer := container.NewHBox(sessionlabel, sessioninput)
	switchbutton := widget.NewButton("Switch!", func() {
		if userinput.Text != "" {
			user = userinput.Text
		}
		if sessioninput.Text != "" {
			session = sessioninput.Text
		}
		cli.SwitchSession(&pb.ClientNeededSession{
			User:    user,
			Session: session,
		})
	})
	return container.NewVBox(usercontainer, sessioncontainer, switchbutton)
}

func ShowNodes(w fyne.Window, cli Client, objects ...fyne.CanvasObject) {
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
		w.SetContent(container.NewVBox(append([]fyne.CanvasObject{list}, objects...)...))
		<-time.After(time.Second)
	}
}
