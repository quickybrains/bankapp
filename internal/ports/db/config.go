package db

import (
	"fmt"
	"time"
)

const (
	maxOpenConnections int   = 10
	poolConnectTimeout int   = 3
	queryTimeout       int64 = 3
)

type Config struct {
	Path     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`

	MaxOpenConnections       int    `yaml:"maxOpenConn"`
	PoolConnectTimeoutInSecs int    `yaml:"poolConnectTimeout"`
	SSLMode                  string `yaml:"sslMode"`

	QueryTimeoutInSecs int64 `yaml:"queryTimeout"`
}

func (c *Config) DSN() string {
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Path,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

// TODO: add options
func (c *Config) DSNWithOptions() string {
	if c.MaxOpenConnections <= 0 {
		c.MaxOpenConnections = maxOpenConnections
	}

	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?pool_max_conns=%v&sslmode=%s",
		c.User,
		c.Password,
		c.Path,
		c.Port,
		c.DBName,
		c.MaxOpenConnections,
		c.SSLMode,
	)
}

func (c *Config) PoolConnectTimeoutInSeconds() time.Duration {
	if c.PoolConnectTimeoutInSecs < poolConnectTimeout {
		return time.Duration(poolConnectTimeout) * time.Second
	}

	return time.Duration(c.PoolConnectTimeoutInSecs) * time.Second
}

func (c *Config) QueryTimeoutInSeconds() time.Duration {
	if c.QueryTimeoutInSecs < queryTimeout {
		return time.Duration(queryTimeout) * time.Second
	}

	return time.Duration(c.QueryTimeoutInSecs) * time.Second
}
