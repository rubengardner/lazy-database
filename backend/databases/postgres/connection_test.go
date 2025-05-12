// postgres_test.go
package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestConfig() *PostgresConfig {
	return &PostgresConfig{
		User:     "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     5432,
		DBName:   "dbhvu01u3e3fug",
		SSLMode:  "disable",
	}
}

func TestGetTableStructureRealDB(t *testing.T) {
	// Setup
	config := getTestConfig()
	conn, err := NewDatabaseConnection(config)
	assert.NoError(t, err)
	defer conn.Db.Close()

	columns, err := conn.GetTableStructure("users")
	assert.NoError(t, err)
	assert.Greater(t, len(columns), 0)
}

func TestGetTableDataRealDB(t *testing.T) {
	// Setup
	config := getTestConfig()
	conn, err := NewDatabaseConnection(config)
	assert.NoError(t, err)
	defer conn.Db.Close()

	// Assume table "users" exists and has some data
	data, err := conn.GetTableData("users")
	assert.NoError(t, err)

	// You can assert specific content here depending on your setup
	t.Logf("Fetched %d rows from users table", len(data))
}
