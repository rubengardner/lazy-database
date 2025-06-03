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
	tableData     [][]string
}

func newModel() model {
	return model{
		onCursor:      0,
		selected:      0,
		connections:   []string{},
		configuration: map[string]*postgres.PostgresConfig{},
		tablesCursor:  0,
		tables:        []string{},
		tableData:     [][]string{},
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

	if err := layout(g, &m, db_connection); err != nil {
		log.Fatal(err)
	}

	if err := keybindings(g, &m, db_connection); err != nil {
		log.Fatal(err)
	}

	if _, err := g.SetCurrentView("Connections"); err != nil {
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
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		// Render the connections list
		updateConnectionsView(v, m)
	}
	if v, err := g.SetView("Tables", maxX/4, 0, maxX/2, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Database Tables"
		v.Autoscroll = true
		v.Wrap = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		updateTablesView(v, m, connection)
	}
	if v, err := g.SetView("Data", maxX/2, 0, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Table Data"
		v.Autoscroll = true
		v.Wrap = true

		updateDataView(v, m)
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
	if err != nil {
		fmt.Println(err)
	}
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

func updateDataView(v *gocui.View, m *model) {
	v.Clear()

	if len(m.tableData) == 0 {
		fmt.Fprintf(v, "Select a table to view data")
		return
	}

	for _, row := range m.tableData {
		fmt.Fprintf(v, "%s\n", formatRow(row))
	}
}

func formatRow(row []string) string {
	// Simple formatter - can be enhanced for better display
	return fmt.Sprintf("%v", row)
}

func keybindings(g *gocui.Gui, m *model, connection *postgres.DatabaseConnection) error {
	// Quit binding
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	// Navigation in Connections view
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

	// Connection selection
	if err := g.SetKeybinding("Connections", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.connections) > 0 {
			m.selected = m.onCursor
			tables, err := connection.GetAllTables()
			if err == nil {
				m.tables = tables
				m.tablesCursor = 0
			}
			updateViews(g, m, connection)
			if _, err := g.SetCurrentView("Tables"); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Connections", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Tables"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Tables", gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Connections"); err != nil {
			return err
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

	if err := g.SetKeybinding("Tables", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.tables) > 0 {
			tableName := m.tables[m.tablesCursor]
			tableData, err := connection.GetTableData(tableName)
			if err == nil {
				m.tableData = [][]string{}
				m.tableData = [][]string{}

				if tableData != nil {
					headers := tableData.Headers
					if len(headers) > 0 {
						m.tableData = append(m.tableData, headers)
					}

					for _, row := range tableData.Rows {
						rowData := []string{}
						for _, val := range row {
							rowData = append(rowData, fmt.Sprintf("%v", val))
						}
						m.tableData = append(m.tableData, rowData)
					}
				}
			}
			updateViews(g, m, connection)
			if _, err := g.SetCurrentView("Data"); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Data", gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Tables"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Tables", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Data"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func updateViews(g *gocui.Gui, m *model, connection *postgres.DatabaseConnection) {
	if v, err := g.View("Connections"); err == nil {
		updateConnectionsView(v, m)
	}
	if v, err := g.View("Tables"); err == nil {
		updateTablesView(v, m, connection)
	}
	if v, err := g.View("Data"); err == nil {
		updateDataView(v, m)
	}
}
