package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/quickybrains/bankapp/internal/app"
)

func main() {
	confPath := flag.String("c", "./config/config.yaml", "path to config file")

	flag.Parse()

	fmt.Println("app started with config path", *confPath)

	config, err := app.NewConfig(*confPath)
	if err != nil {
		panic(err)
	}

	node, err := app.StartNewNode(config)
	if err != nil {
		log.Fatalf("app start new node, err=%v", err)
	}

	// Stopping the node on these signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	node.Stop()
}
