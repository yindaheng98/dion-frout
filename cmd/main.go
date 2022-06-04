package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"github.com/yindaheng98/dion-frout/pkg/controller"
	client "github.com/yindaheng98/dion-frout/pkg/dion"
	"github.com/yindaheng98/dion/pkg/islb"
)

func main() {
	a := app.New()
	cli := &client.Client{Node: islb.NewNode("test")}
	go controller.Control(a, cli)
	a.Run()
	fmt.Println("Exiting")
}
