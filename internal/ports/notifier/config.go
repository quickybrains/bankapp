package notifier

import "fmt"

const (
	defaultHost     = "localhost"
	defaultPort     = 5672
	defaultUser     = "guest"
	defaultPassword = "guest"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (c *Config) DSN() string {
	if c.Host == "" {
		c.Host = defaultHost
	}

	if c.Port == 0 {
		c.Port = defaultPort
	}

	if c.User == "" {
		c.User = defaultUser
	}

	if c.Password == "" {
		c.Password = defaultPassword
	}

	return fmt.Sprintf(
		"amqp://%s:%s@%s:%v/",
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}
