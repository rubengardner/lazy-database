package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabaseConnection struct {
	Db *sql.DB
}

func NewDatabaseConnection(config *PostgresConfig) (*DatabaseConnection, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &DatabaseConnection{Db: db}, nil
}

func (db *DatabaseConnection) GetTableStructure(tableName string) ([]string, error) {
	rows, err := db.Db.Query(`
  SELECT column_name
  FROM information_schema.columns
  WHERE table_name = $1 AND table_schema = 'public'
`, tableName)
	if err != nil {
		return nil, err
	}
	var columns []string
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, err
		}
		columns = append(columns, columnName)
	}

	return columns, nil
}

func (db *DatabaseConnection) GetAllTables() ([]string, error) {
	rows, err := db.Db.Query(`
  SELECT table_name
  FROM information_schema.tables
  WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
  ORDER BY table_name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}

type TableData struct {
	Headers []string        // Column names
	Rows    [][]interface{} // Row data as array of values
}

// GetTableData fetches all data from a given table
func (db *DatabaseConnection) GetTableData(tableName string) (*TableData, error) {
	rows, err := db.Db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get the column names from the result
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := &TableData{
		Headers: columns,
		Rows:    [][]interface{}{},
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range columns {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return nil, err
		}

		// Map column names to the values
		rowData := make(map[string]any)
		for i, columnName := range columns {
			val := values[i]
			rowData[columnName] = val
		}

		result.Rows = append(result.Rows, values)
	}

	return result, nil
}
