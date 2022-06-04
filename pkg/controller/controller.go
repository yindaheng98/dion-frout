package controller

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

func Control(a fyne.App, cli Client) {
	addr := GetNatsAddr(a)
	log.Println("Connecting: ", addr)
	if err := cli.Connect(addr); err != nil {
		log.Fatalln(err)
	}

	w := a.NewWindow("dion system")
	session := SessionEntry(cli)
	button := widget.NewButton("Refresh!", func() {})
	refresh := func() {
		w.SetContent(container.NewGridWithRows(2, container.NewVBox(session, button), RefreshNodeTable(cli)))
	}
	button.OnTapped = refresh
	refresh()
	w.Resize(fyne.NewSize(800, 600))
	w.Show()
}

func GetNatsAddr(a fyne.App) string {
	w := a.NewWindow("dion system")
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
	w.Resize(fyne.NewSize(200, 100))
	w.Show()
	addr = <-addrCh
	w.Hide()
	w.Close()
	return addr
}

func SessionEntry(cli Client) *fyne.Container {
	user := "stupid"
	userlabel := widget.NewLabel("User:")
	userinput := widget.NewEntry()
	userinput.SetPlaceHolder(user)
	session := "stupid"
	sessionlabel := widget.NewLabel("Session:")
	sessioninput := widget.NewEntry()
	sessioninput.SetPlaceHolder(session)
	inputcontainer := container.New(layout.NewFormLayout(), userlabel, userinput, sessionlabel, sessioninput)
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
	box := container.NewVBox(inputcontainer, switchbutton)
	return box
}

func RefreshNodeTable(cli Client) *widget.Table {
	log.Println("Refreshing node table")
	head := []string{
		"NID", "Service", "DC", "ExtraInfo",
	}
	var ids []string
	var nodes [][]string
	for id, node := range cli.GetNodes() {
		ids = append(ids, id)
		nodes = append(nodes, []string{
			node.NID, node.Service, node.DC,
			fmt.Sprintf("%+v", node.ExtraInfo),
		})
	}
	list := widget.NewTable(
		func() (int, int) {
			return len(nodes) + 1, len(head)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("unknown unknown", func() {})
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			button := o.(*widget.Button)
			if i.Row == 0 {
				button.SetText(head[i.Col])
				button.OnTapped = func() {
					fmt.Println("This is head")
				}
			} else {
				button.SetText(nodes[i.Row-1][i.Col])
				button.OnTapped = func() {
					cli.SwitchNode(ids[i.Row-1])
				}
			}
		})
	return list
}
