package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

type model struct {
	connections map[int]string
	selected    int
}

// Create a new instance of the model
func newModel() model {
	return model{
		connections: map[int]string{1: "Database1", 2: "Database2", 3: "Database3"},
		selected:    0,
	}
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	m := newModel()

	g.SetManagerFunc(func(g *gocui.Gui) error {
		return layout(g, &m)
	})

	if err := keybindings(g, &m); err != nil {
		log.Fatal(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatal(err)
	}
}

func layout(g *gocui.Gui, m *model) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("Connections", 0, 0, maxX/4, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Databases"
		v.Autoscroll = false
		v.Wrap = true

		// Render the connections list
		updateConnectionsView(v, m)
	}

	if v, err := g.SetView("right", 31, 0, 80, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Database Details"
		v.Autoscroll = true
		v.Wrap = true

		updateDetailsView(v, m)
	}

	return nil
}

func updateConnectionsView(v *gocui.View, m *model) {
	v.Clear()
	for i, db := range m.connections {
		cursor := " "
		if i == m.selected {
			cursor = ">"
		}
		fmt.Fprintf(v, "%s %s\n", cursor, db)
	}
}

func updateDetailsView(v *gocui.View, m *model) {
	v.Clear()
	if len(m.connections) > 0 {
		fmt.Fprintln(v, "Selected database: ", m.connections[m.selected])
		fmt.Fprintln(v, m.selected)
	}
}

func keybindings(g *gocui.Gui, m *model) error {
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.selected < len(m.connections)-1 {
			m.selected++
			updateViews(g, m)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.selected > 0 {
			m.selected--
			updateViews(g, m)
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.connections) > 0 {
			m.selected = m.selected % len(m.connections)
			updateViews(g, m)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func updateViews(g *gocui.Gui, m *model) {
	if v, err := g.View("Connections"); err == nil {
		updateConnectionsView(v, m)
	}
	if v, err := g.View("right"); err == nil {
		updateDetailsView(v, m)
	}
}
