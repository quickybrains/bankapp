package db

import (
	"net/url"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDSN(t *testing.T) {
	config := &Config{
		Path:     "localhost",
		Port:     "5432",
		User:     "user@host",
		Password: "p@ss/word",
		DBName:   "bank",
	}

	dsn := config.DSN()
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse DSN: %v", err)
	}
	if parsed.Host != "localhost:5432" || parsed.Path != "/bank" {
		t.Fatalf("unexpected DSN target: %s", dsn)
	}
	username := parsed.User.Username()
	password, _ := parsed.User.Password()
	if username != "user@host" || password != "p@ss/word" {
		t.Fatal("unexpected DSN credentials")
	}
	if parsed.Query().Get("sslmode") != "disable" {
		t.Fatalf("expected sslmode=disable, got %q", parsed.Query().Get("sslmode"))
	}
	if _, err := pgxpool.ParseConfig(dsn); err != nil {
		t.Fatalf("pgxpool.ParseConfig: %v", err)
	}
}
