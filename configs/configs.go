package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`

	Logger struct {
		Mode string `json:"mode"`
	} `json:"logger"`

	Postgres struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"postgres"`

	Tarantool struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"tarantool"`
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

func (c *Config) GetServerHost() string {
	return c.Server.Host
}

func (c *Config) GetServerPort() string {
	return c.Server.Port
}

func (c *Config) GetLoggerMode() string {
	return c.Logger.Mode
}

func (c *Config) GetPostgresConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.Name)
}

func (c *Config) GetTarantoolUser() string {
	return c.Tarantool.User
}

func (c *Config) GetTarantoolPassword() string {
	return c.Tarantool.Password
}

func (c *Config) GetTarantoolConfig() string {
	return fmt.Sprint(c.Tarantool.Host, ":", c.Tarantool.Port)
}
