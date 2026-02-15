package http

import (
	"fmt"
	"time"
)

const (
	defaultPath             string = ""
	defaultPort             int    = 8080
	defaultReadTimeoutInMs  int    = 1000
	defaultWriteTimeoutInMs int    = 1000
	defaultReqBodyLimit     int64  = 5126
)

type Config struct {
	path             string `yaml:"host"`
	port             int    `yaml:"port"`
	readTimeoutInMs  int    `yaml:"readTimeout"`
	writeTimeoutInMs int    `yaml:"writeTimeout"`
	reqBodyLimit     int64  `yaml:"bodyLimit"`
}

func (c *Config) ServerURL() string {
	if c.path == "" {
		c.path = defaultPath
	}

	if c.port == 0 {
		c.port = defaultPort
	}

	return fmt.Sprintf("%s:%v", c.path, c.port)
}

func (c *Config) ReadTimeout() time.Duration {
	if c.readTimeoutInMs == 0 {
		c.readTimeoutInMs = defaultReadTimeoutInMs
	}

	return time.Duration(c.readTimeoutInMs) * time.Millisecond
}

func (c *Config) WriteTimeout() time.Duration {
	if c.writeTimeoutInMs == 0 {
		c.writeTimeoutInMs = defaultWriteTimeoutInMs
	}

	return time.Duration(c.writeTimeoutInMs) * time.Millisecond
}

func (c *Config) HTTPReqBodyLimit() int64 {
	if c.reqBodyLimit < defaultReqBodyLimit {
		return defaultReqBodyLimit
	}

	return c.reqBodyLimit
}
