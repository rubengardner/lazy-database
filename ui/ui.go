package main

import (
	"fmt"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

type model struct {
	onCursor      int
	selected      int
	connections   []string
	configuration map[string]*postgres.PostgresConfig
}

func newModel() model {
	return model{
		onCursor:      0,
		selected:      0,
		connections:   []string{},
		configuration: map[string]*postgres.PostgresConfig{},
	}
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	m := newModel()
	loadConfig(&m)

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

	if v, err := g.SetView("Tables", 31, 0, 80, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Database Tables"
		v.Autoscroll = true
		v.Wrap = true

		updateTablesView(v, m)
	}

	return nil
}

func loadConfig(m *model) {
	fileName := "config.json"
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading the file: %v", err)
	}
	configuration, err := postgres.LoadConfig(data)
	m.configuration["postgres"] = configuration
	m.connections = append(m.connections, configuration.Name)
}

func updateConnectionsView(v *gocui.View, m *model) {
	v.Clear()
	for i, db := range m.connections {
		cursor := " "
		if i == m.onCursor {
			cursor = ">"
		}
		fmt.Fprintf(v, "%s %s\n", cursor, db)
	}
}

func updateTablesView(v *gocui.View, m *model) {
	v.Clear()
	if len(m.connections) > 0 {
		fmt.Fprintln(v, "Selected database: ", m.connections[m.selected])
		fmt.Fprintln(v, m.selected)
	}
}

func keybindings(g *gocui.Gui, m *model) error {
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.onCursor < len(m.connections)-1 {
			m.onCursor++
			updateViews(g, m)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.onCursor > 0 {
			m.onCursor--
			updateViews(g, m)
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.connections) > 0 {
			m.selected = m.onCursor
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
	if v, err := g.View("Tables"); err == nil {
		updateTablesView(v, m)
	}
}
