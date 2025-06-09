package ui

import (
	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
)

func Layout(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
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
		v.Autoscroll = false
		v.Wrap = false
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		updateDataView(v, m)
	}
	return nil
}
