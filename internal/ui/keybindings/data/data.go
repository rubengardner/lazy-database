package data

import (
	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
)

func DataKeybindings(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
	if err := MovementDataKeybindings(g, m, connection); err != nil {
		return err
	}
	if err := QueryKeybindings(g, m, connection); err != nil {
		return err
	}
	if err := YankKeybinding(g, m, connection); err != nil {
		return err
	}
	if err := UpdateValueDataKeybindings(g, m, connection); err != nil {
		return err
	}
	return nil
}
