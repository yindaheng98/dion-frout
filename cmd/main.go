package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"github.com/yindaheng98/dion-frout/pkg/controller"
)

func main() {
	a := app.New()
	fmt.Println(controller.GetNatsAddr(a))
}
