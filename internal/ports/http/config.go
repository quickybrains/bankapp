package http

import (
	"fmt"
	"time"
)

const (
	defaultPath             string = "localhost"
	defaultPort             int    = 8080
	defaultReadTimeoutInMs  int    = 1000
	defaultWriteTimeoutInMs int    = 1000
	defaultReqBodyLimit     int64  = 5126
)

type Config struct {
	Path             string `yaml:"host"`
	Port             int    `yaml:"port"`
	ReadTimeoutInMs  int    `yaml:"readTimeout"`
	WriteTimeoutInMs int    `yaml:"writeTimeout"`
	ReqBodyLimit     int64  `yaml:"bodyLimit"`
}

func (c *Config) ServerURL() string {
	if c.Path == "" {
		c.Path = defaultPath
	}

	if c.Port == 0 {
		c.Port = defaultPort
	}

	return fmt.Sprintf("%s:%v", c.Path, c.Port)
}

func (c *Config) ReadTimeout() time.Duration {
	if c.ReadTimeoutInMs == 0 {
		c.ReadTimeoutInMs = defaultReadTimeoutInMs
	}

	return time.Duration(c.ReadTimeoutInMs) * time.Millisecond
}

func (c *Config) WriteTimeout() time.Duration {
	if c.WriteTimeoutInMs == 0 {
		c.WriteTimeoutInMs = defaultWriteTimeoutInMs
	}

	return time.Duration(c.WriteTimeoutInMs) * time.Millisecond
}

func (c *Config) HTTPReqBodyLimit() int64 {
	if c.ReqBodyLimit < defaultReqBodyLimit {
		return defaultReqBodyLimit
	}

	return c.ReqBodyLimit
}
