package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/quickybrains/bankapp/internal/app"
)

var (
	config     = flag.String("c", "./config/config.yaml", "path to config file")
	migrations = flag.String("m", "./internal/ports/db/migrations", "path migrations dir")
)

func main() {
	flag.Parse()

	conf, err := app.NewConfig(*config)
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New(fmt.Sprintf("file://%s", *migrations), conf.DB.DSN())
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 1 {
		log.Fatal("Migration command is not specified (up/down)")
	}

	switch {
	case os.Args[1] == "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case os.Args[1] == "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown migration command %s. Only up/down available.", os.Args[1])
	}
}
