package main

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ikhwan11/auth-be/internal/config"
	"github.com/ikhwan11/auth-be/internal/database"
)

func buildAuthDSN(cfg config.AppConfig) string {
	db := cfg.AuthDB
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=5",
		db.User, db.Password, db.Host, db.Port, db.DBName, db.SSLMode,
	)
}

func runMigration(direction string) error {
	cfg := config.Load()

	if db, err := database.NewConnection(cfg.AuthDB); err != nil {
		return fmt.Errorf("failed to connect auth db before migration: %w", err)
	} else {
		defer db.Close()
	}

	dsn := buildAuthDSN(cfg)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("failed to init migrate: %w", err)
	}

	switch direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		return fmt.Errorf("unknown direction: %s", direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("Migration", direction, "berhasil dijalankan ke auth db")
	return nil
}
