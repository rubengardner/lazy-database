package postgres

import (
	"encoding/json"
	"fmt"
)

type PostgresConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

func LoadConfig(data []byte) (*PostgresConfig, error) {
	var config PostgresConfig
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}
	return &config, err
}
