package db

import (
	"net"
	"net/url"
	"strconv"
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
	return c.connectionURL(url.Values{})
}

func (c *Config) DSNWithOptions() string {
	maxConnections := c.MaxOpenConnections
	if maxConnections <= 0 {
		maxConnections = maxOpenConnections
	}

	query := url.Values{
		"pool_max_conns": []string{strconv.Itoa(maxConnections)},
	}
	return c.connectionURL(query)
}

func (c *Config) connectionURL(query url.Values) string {
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	query.Set("sslmode", sslMode)

	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Path, c.Port),
		Path:     "/" + c.DBName,
		RawQuery: query.Encode(),
	}).String()
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
