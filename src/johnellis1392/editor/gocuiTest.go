package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func onPressEnter(g *gocui.Gui, v *gocui.View) error {
	fmt.Fprintln(v, "Enter Clicked!")
	return nil
}

// GoCUITest is an example of using gocui package.
func GoCUITest() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	if v, err := g.SetView("viewname", 2, 2, 22, 7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "This is a New View")
	}

	if err := g.SetKeybinding("viewname", gocui.KeyEnter, gocui.ModNone, onPressEnter); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
