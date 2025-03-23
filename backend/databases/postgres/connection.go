package postgres

import (
	"database/sql"
	"fmt"
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
	rows, err := db.Db.Query(fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s'", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// GetTableData fetches all data from a given table
func (db *DatabaseConnection) GetTableData(tableName string) ([]map[string]interface{}, error) {
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

	var results []map[string]interface{}
	for rows.Next() {
		// Create a slice of interface{} to hold the row values
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range columns {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return nil, err
		}

		// Map column names to the values
		rowData := make(map[string]interface{})
		for i, columnName := range columns {
			val := values[i]
			rowData[columnName] = val
		}

		results = append(results, rowData)
	}

	return results, nil
}
