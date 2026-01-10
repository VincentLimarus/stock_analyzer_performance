package migrations

import (
	configs "VincentLimarus/stock-analyzer-performance/config"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run execute when u need to use auto Migration for database
// Recommended for development and staging
func Run(ctx context.Context) error {
	// init database migration
	source := "file://migrations"

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		configs.Env.DBUser,
		configs.Env.DBPassword,
		configs.Env.DBHost,
		configs.Env.DBPort,
		configs.Env.DBName,
	)

	log.Println("Migration DSN:", dsn)

	m, err := migrate.New(source, dsn)
	if err != nil {
		log.Printf("error init golang-migrate: %v", err)
		return err
	}
	defer m.Close()

	// if migration has no change, ignore error no change
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migration error: %v", err)
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}
