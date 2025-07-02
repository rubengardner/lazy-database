package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
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
	Headers []string
	Rows    [][]interface{}
}

func (db *DatabaseConnection) GetTableData(tableName string, whereClause ...string) (*TableData, error) {
	query := fmt.Sprintf("SELECT * FROM %s", pq.QuoteIdentifier(tableName))
	if len(whereClause) > 0 && whereClause[0] != "" {
		query = fmt.Sprintf("%s WHERE %s", query, whereClause[0])
	}
	query = fmt.Sprintf("%s LIMIT 50", query)
	rows, err := db.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
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

		row := make([]interface{}, len(columns))
		for i, val := range values {
			if val != nil {
				if bytes, ok := val.([]byte); ok {
					colType := columnTypes[i].DatabaseTypeName()
					if colType == "JSON" || colType == "JSONB" {
						row[i] = string(bytes)
					} else {
						row[i] = string(bytes)
					}
				} else {
					row[i] = val
				}
			} else {
				row[i] = val
			}
		}
		result.Rows = append(result.Rows, row)
	}
	return result, nil
}

func (db *DatabaseConnection) GetPrimaryKeyColumn(tableName string) (string, error) {
	query := `
	SELECT a.attname
	FROM   pg_index i
	JOIN   pg_attribute a ON a.attrelid = i.indrelid
								AND a.attnum = ANY(i.indkey)
	WHERE  i.indrelid = $1::regclass
	AND    i.indisprimary;
	`

	row := db.Db.QueryRow(query, tableName)

	var primaryKey string
	err := row.Scan(&primaryKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("table %s has no primary key", tableName)
		}
		return "", err
	}

	return primaryKey, nil
}

func (db *DatabaseConnection) UpdateTableValue(tableName, columnName, primaryKeyColumn string, primaryKeyValue interface{}, newValue interface{}) error {
	query := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE %s = $2",
		pq.QuoteIdentifier(tableName),
		pq.QuoteIdentifier(columnName),
		pq.QuoteIdentifier(primaryKeyColumn))

	_, err := db.Db.Exec(query, newValue, primaryKeyValue)
	if err != nil {
		return fmt.Errorf("failed to update value: %v", err)
	}

	return nil
}
