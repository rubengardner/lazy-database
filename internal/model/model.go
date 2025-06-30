package model

import "github.com/rubengardner/lazy-database/backend/databases/postgres"

type LazyDBState struct {
	OnCursor      int
	Selected      int
	Connections   []string
	Configuration map[string]*postgres.PostgresConfig
	TablesCursor  int
	Tables        []string
	TableData     [][]string
	DataCursorRow int // Current row position in data view
	DataCursorCol int // Current column position in data view
}

func NewLazyDBState() LazyDBState {
	return LazyDBState{
		OnCursor:      0,
		Selected:      0,
		Connections:   []string{},
		Configuration: map[string]*postgres.PostgresConfig{},
		TablesCursor:  0,
		Tables:        []string{},
		TableData:     [][]string{},
		DataCursorRow: 1, // Start at row 1 (below header)
		DataCursorCol: 0,
	}
}
