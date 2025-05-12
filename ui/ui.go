package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rubengardner/lazy-database/backend/databases/postgres"

	"github.com/jroimartin/gocui"
)

type model struct {
	onCursor      int
	selected      int
	connections   []string
	configuration map[string]*postgres.PostgresConfig
	tablesCursor  int
	tables        []string
}

func newModel() model {
	return model{
		onCursor:      0,
		selected:      0,
		connections:   []string{},
		configuration: map[string]*postgres.PostgresConfig{},
		tablesCursor:  0,
		tables:        []string{},
	}
}

func main() {
	m := newModel()
	loadConfig(&m)
	db_connection, err := postgres.NewDatabaseConnection(m.configuration["postgres"])
	if err != nil {
		log.Fatal(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	g.SetManagerFunc(func(g *gocui.Gui) error {
		return layout(g, &m, db_connection)
	})
	if err := keybindings(g, &m, db_connection); err != nil {
		log.Fatal(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatal(err)
	}
}

func layout(g *gocui.Gui, m *model, connection *postgres.DatabaseConnection) error {
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

	if v, err := g.SetView("Tables", maxX/4, 0, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Database Tables"
		v.Autoscroll = true
		v.Wrap = true

		updateTablesView(v, m, connection)
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

func updateTablesView(v *gocui.View, m *model, connection *postgres.DatabaseConnection) {
	v.Clear()
	if len(m.connections) == 0 {
		fmt.Fprintf(v, "No database connections available")
		return
	}

	if len(m.tables) == 0 {
		tables, err := connection.GetAllTables()
		if err != nil {
			fmt.Fprintf(v, "Error fetching tables: %v\n", err)
			return
		}
		m.tables = tables
	}

	for i, table := range m.tables {
		cursor := " "
		if i == m.tablesCursor {
			cursor = ">"
		}
		fmt.Fprintf(v, "%s %s\n", cursor, table)
	}
}

func keybindings(g *gocui.Gui, m *model, connection *postgres.DatabaseConnection) error {
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.onCursor < len(m.connections)-1 {
			m.onCursor++
			updateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.onCursor > 0 {
			m.onCursor--
			updateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.connections) > 0 {
			m.selected = m.onCursor
			if _, err := g.View("Tables"); err == nil {
				tables, err := connection.GetAllTables()
				if err == nil {
					m.tables = tables
					m.tablesCursor = 0
				}
			}
			updateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Tables"); err == nil {
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Connections"); err == nil {
			// Focus on connections view
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Tables", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.tablesCursor < len(m.tables)-1 {
			m.tablesCursor++
			updateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Tables", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.tablesCursor > 0 {
			m.tablesCursor--
			updateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}
	// if err := g.SetKeybinding("Tables", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
	// 		m.table[m.tablesCursor]
	// 		updateViews(g, m, connection)
	// 	}
	// 	return nil
	// }); err != nil {
	// 	return err
	// }

	return nil
}

func updateViews(g *gocui.Gui, m *model, connection *postgres.DatabaseConnection) {
	if v, err := g.View("Connections"); err == nil {
		updateConnectionsView(v, m)
	}
	if v, err := g.View("Tables"); err == nil {
		updateTablesView(v, m, connection)
	}
}
