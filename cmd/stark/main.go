package main

import (
	"log"

	"github.com/hackfeed/stark/internal/client"
	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalln(err)
	}
	defer g.Close()

	g.SetManagerFunc(client.Layout)
	g.SetKeybinding("name", gocui.KeyEnter, gocui.ModNone, client.Connect)
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, client.Send)
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, client.Disconnect)
	g.MainLoop()
}
