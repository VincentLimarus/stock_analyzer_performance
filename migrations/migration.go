package migrations

import (
	configs "VincentLimarus/stock-analyzer-performance/config"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
)

// Run execute when u need to use auto Migration for database
// Recommended for development and staging
func Run(ctx context.Context) error {
	const logCtx = "migrations.migration.Run"

	// init database migration
	// use DB master to migrate.
	source := "file://migrations"
	host := configs.Env.DBHost
	port := configs.Env.DBPort
	user := configs.Env.DBUser
	password := configs.Env.DBPassword
	dbname := configs.Env.DBName

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port)

	m, err := migrate.New(source, dsn)
	if err != nil {
		log.Printf("error init golang-migrate: %v", err)
		return err
	}

	// if migration has no change, ignore error no change
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
