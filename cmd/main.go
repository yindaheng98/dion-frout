package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"github.com/yindaheng98/dion-frout/pkg/controller"
	client "github.com/yindaheng98/dion-frout/pkg/dion"
)

func main() {
	a := app.New()
	go controller.Control(a, client.NewClient("test"))
	a.Run()
	fmt.Println("Exiting")
}
