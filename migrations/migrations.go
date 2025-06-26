package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

// Initialize configura el sistema de migraciones
// Updated to include RBAC migration
func Initialize(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Debug: List embedded migrations
	entries, err := embedMigrations.ReadDir(".")
	if err == nil {
		fmt.Printf("Embedded migrations found: %d\n", len(entries))
		for _, entry := range entries {
			fmt.Printf("- %s\n", entry.Name())
		}
	}

	return nil
}

// Up ejecuta todas las migraciones pendientes
func Up(db *sql.DB) error {
	if err := Initialize(db); err != nil {
		return err
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Down revierte todas las migraciones
func Down(db *sql.DB) error {
	if err := Initialize(db); err != nil {
		return err
	}

	if err := goose.Down(db, "."); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}
