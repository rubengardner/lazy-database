package ui

import (
	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
)

func Layout(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
	maxX, maxY := g.Size()
	currentView := g.CurrentView()
	isTablesViewSelected := currentView != nil && currentView.Name() == "Tables"
	isDataViewSelected := currentView != nil && currentView.Name() == "Data"
	isDetailView := isTablesViewSelected || isDataViewSelected || m.InDetailView

	m.InDetailView = isDetailView
	if !isDetailView {
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
		// Tables view takes 1/4 of screen
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
	} else {
		m.InDetailView = true
		g.DeleteView("Connections")

		// Tables view now takes 1/4 of screen at the left edge
		if v, err := g.SetView("Tables", 0, 0, maxX/4, maxY-1); err != nil {
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
	}

	dataStartX := maxX / 4
	if !isDetailView {
		dataStartX = maxX / 2
	}
	if v, err := g.SetView("Data", dataStartX, 0, maxX, maxY-1); err != nil {
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
