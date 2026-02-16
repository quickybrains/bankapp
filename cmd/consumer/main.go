package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/quickybrains/bankapp/internal/app"
	"github.com/quickybrains/bankapp/internal/ports/log/zerolog"
	"github.com/quickybrains/bankapp/internal/ports/notifier/rabbitmq"
)

func main() {
	confPath := flag.String("c", "./config/config_local.yaml", "path to config file")

	flag.Parse()

	fmt.Println("app started with config path", *confPath)

	config, err := app.NewConfig(*confPath)
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(config.Log).With("service", "consumer")

	c, err := rabbitmq.StartNewConsumer(config.Notifier, logger)

	// Stopping the node on these signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	if err := c.Stop(); err != nil {
		panic(err)
	}
}
