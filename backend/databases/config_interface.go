package config

type DatabaseConfig interface {
	LoadConfig(filePath string) (interface{}, error)
}
