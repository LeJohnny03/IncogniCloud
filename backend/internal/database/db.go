package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	// NEU: Importiere golang-migrate und die benötigten Treiber
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Erfolgreich mit PostgreSQL verbunden!")
	return db, nil
}

// NEU: Diese Funktion führt die Migrationen aus
func RunMigrations(db *sqlx.DB, dbName string) error {
	// Erstelle den Datenbank-Treiber für die Migration
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("fehler beim Erstellen des Migration-Treibers: %w", err)
	}

	// Zeige auf den Ordner mit deinen .sql Dateien
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Initialisieren der Migration: %w", err)
	}

	// Führe alle "up" Migrationen aus
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("fehler beim Ausführen der Migration: %w", err)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("Datenbank ist bereits auf dem neuesten Stand (Keine neuen Migrationen).")
	} else {
		fmt.Println("Migrationen erfolgreich ausgeführt! Tabellen wurden angelegt/aktualisiert.")
	}

	return nil
}
