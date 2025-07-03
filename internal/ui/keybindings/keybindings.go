package keybindings

import (
	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
	"github.com/rubengardner/lazy-database/internal/ui/keybindings/connections"
	"github.com/rubengardner/lazy-database/internal/ui/keybindings/core"
	"github.com/rubengardner/lazy-database/internal/ui/keybindings/data"
	"github.com/rubengardner/lazy-database/internal/ui/keybindings/tables"
)

func Keybindings(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
	if err := core.CoreKeybindings(g); err != nil {
		return err
	}
	if err := connections.MovementConnectionsKeybindings(g, m, connection); err != nil {
		return err
	}
	if err := tables.MovementTablesKeybindings(g, m, connection); err != nil {
		return err
	}
	if err := data.DataKeybindings(g, m, connection); err != nil {
		return err
	}
	return nil
}
