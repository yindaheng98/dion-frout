package controller

import (
	"context"
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
	Connect(nats_addr, ffplay_path string) error
	GetNodes() map[string]discovery.Node
	SwitchNode(id string)
	SwitchSession(session *pb.ClientNeededSession)
}

func Control(a fyne.App, cli Client) {
	addr, path := Init(a)
	log.Println("Connecting: ", addr)
	if err := cli.Connect(addr, path); err != nil {
		log.Fatalf("Cannot Connect: %+v\n", err)
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

func Init(a fyne.App) (string, string) {
	w := a.NewWindow("dion system")
	addr := "nats://127.0.0.1:4222"
	addrlabel := widget.NewLabel("NATS Address:")
	addrinput := widget.NewEntry()
	addrinput.SetPlaceHolder(addr)
	path := "ffplay"
	pathlabel := widget.NewLabel("FFplay command:")
	pathinput := widget.NewEntry()
	pathinput.SetPlaceHolder(path)
	ctx, cancel := context.WithCancel(context.Background())
	connect := widget.NewButton("Connect!", func() {
		if addrinput.Text != "" {
			addr = addrinput.Text
		}
		if pathinput.Text != "" {
			path = pathinput.Text
		}
		cancel()
	})
	inputcontainer := container.New(layout.NewFormLayout(), addrlabel, addrinput, pathlabel, pathinput)
	w.SetContent(container.NewVBox(inputcontainer, connect))
	w.Resize(fyne.NewSize(400, 100))
	w.Show()
	<-ctx.Done()
	w.Hide()
	w.Close()
	return addr, path
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
					log.Println("This is head")
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
