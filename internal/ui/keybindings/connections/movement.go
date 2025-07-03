package connections

import (
	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
	"github.com/rubengardner/lazy-database/internal/ui"
)

func MovementConnectionsKeybindings(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
	if err := g.SetKeybinding("Connections", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Tables"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Connections", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.Connections) > 0 {
			m.Selected = m.OnCursor
			tables, err := connection.GetAllTables()
			if err == nil {
				m.Tables = tables
				m.TablesCursor = 0
			}
			ui.UpdateViews(g, m, connection)
			if _, err := g.SetCurrentView("Tables"); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Connections", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.OnCursor > 0 {
			m.OnCursor--
			ui.UpdateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
