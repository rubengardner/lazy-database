package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
	"github.com/rubengardner/lazy-database/internal/ui"
)

func main() {
	m := model.NewLazyDBState()
	loadConfig(&m)
	db_connection, err := postgres.NewDatabaseConnection(m.Configuration["postgres"])
	if err != nil {
		log.Fatal(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	g.SetManagerFunc(func(g *gocui.Gui) error {
		return ui.Layout(g, &m, db_connection)
	})

	if err := ui.Layout(g, &m, db_connection); err != nil {
		log.Fatal(err)
	}

	if err := ui.Keybindings(g, &m, db_connection); err != nil {
		log.Fatal(err)
	}

	if _, err := g.SetCurrentView("Connections"); err != nil {
		log.Fatal(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatal(err)
	}
}

func loadConfig(m *model.LazyDBState) {
	fileName := "config.json"
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading the file: %v", err)
	}
	configuration, err := postgres.LoadConfig(data)
	if err != nil {
		fmt.Println(err)
	}
	m.Configuration["postgres"] = configuration
	m.Connections = append(m.Connections, configuration.Name)
}
