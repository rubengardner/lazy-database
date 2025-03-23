package config

import (
	"fmt"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
)

type DATABASES int

const (
	POSTGRES DATABASES = iota
	MYSQL
)

type DatabaseConfigFactory struct{}

func (d *DatabaseConfigFactory) Run(database DATABASES configuration_path string) (DatabaseConfigInterface, error) {

	switch database {
	case POSTGRES:
		config, errors := postgres.LoadConfig(data)
		if errors != nil {
			return nil, errors
		}
		postgres.InitializeDBConnection(config)
	default:
		return nil, fmt.Errorf("unknown database type: %d", database)
	}
}
