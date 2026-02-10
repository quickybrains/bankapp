package app

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/quickybrains/bankapp/internal/ports/db"
	"github.com/quickybrains/bankapp/internal/ports/http"
	"github.com/quickybrains/bankapp/internal/ports/log/zerolog"
	"github.com/quickybrains/bankapp/internal/ports/notifier"
	"gopkg.in/yaml.v2"
)

const shutdownTimeout int64 = 5

type Config struct {
	HTTP     *http.Config     `yaml:"http"`
	DB       *db.Config       `yaml:"db"`
	Log      *zerolog.Config  `yaml:"log"`
	Notifier *notifier.Config `yaml:"notifier"`

	ShutdownTimeoutInSecs int64 `yaml:"shutdownTimeout"`
}

func NewConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os open config file, path=%s", path)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("io read all config file")
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("yaml unmarshal config")
	}

	return config, nil
}

func (c *Config) ShutdownTimeoutInSeconds() time.Duration {
	if c.ShutdownTimeoutInSecs < shutdownTimeout {
		return time.Duration(shutdownTimeout) * time.Second
	}

	return time.Duration(c.ShutdownTimeoutInSecs) * time.Second
}
