package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"database"`

	Logger struct {
		Mode string `json:"mode"`
	} `json:"logger"`

	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
}

func (c *Config) GetDBConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name)
}

func (c *Config) GetLoggerConfig() string {
	return c.Logger.Mode
}

func (c *Config) GetServerHost() string {
	return c.Server.Host
}

func (c *Config) GetServerPort() string {
	return c.Server.Port
}

func LoadConfig(name string) (*Config, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
