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
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%v/",
		defaultUser,
		defaultPassword,
		defaultHost,
		defaultPort,
	)
}
