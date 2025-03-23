package config

type DatabaseConfigInterface interface {
	LoadConfig(data []byte) error
}
