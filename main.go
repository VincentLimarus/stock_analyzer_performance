package main

import (
	"VincentLimarus/go-skeleton-files/client"
	configs "VincentLimarus/go-skeleton-files/config"
	"VincentLimarus/go-skeleton-files/migrations"
	"VincentLimarus/go-skeleton-files/repository"
	"VincentLimarus/go-skeleton-files/service"
	"VincentLimarus/go-skeleton-files/util"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cmdHttp "VincentLimarus/go-skeleton-files/cmd/http"
	delivery "VincentLimarus/go-skeleton-files/delivery/http"
	serviceHealth "VincentLimarus/go-skeleton-files/service/health"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		wg          = sync.WaitGroup{}
	)
	
	// Load environment variables
	configs.Init()
	// End of loading environment variables
	
	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		configs.Env.DBHost, configs.Env.DBUser, configs.Env.DBPassword, configs.Env.DBName, configs.Env.DBPort)

	var err error
	
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	// End of database connection

	// Auto Migration
	if configs.Env.AutoMigration == "true" {
		if err := migrations.Run(context.Background()); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
	}
	// End of Auto Migration

	// Initialize Dependencies
	repositoryRegistry := initRepositoryDependency(repositoryDependency{
		DB: DB,
	})

	serviceRegistry := initServiceDependency(serviceDependency{
		repoRegistry: repositoryRegistry,
		DB:        DB,
	})

	deliveryRegistry := initDeliveryDependency(deliveryDependency{
		serviceRegistry: &serviceRegistry,
	})
	// End of Initialize Dependencies	

	middlewareLimiter := util.CreateRateLimiter(100, time.Minute)
	// Start HTTP Server
	httpServer := cmdHttp.NewServer( 
		deliveryRegistry,
		middlewareLimiter,
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Serve(ctx)

		// Graceful Shutdown //
		shutdownDelay := 10 * time.Second
		if configs.Env.ShutDownDelayInSeconds > 0 {
			shutdownDelay = time.Duration(configs.Env.ShutDownDelayInSeconds) * time.Second
		}

		wait := util.GracefullyShutdown(ctx, shutdownDelay,
			map[string]util.Operation{
				"db": func(_ context.Context) error {
					log.Println("Closing database connection...")
					return DB.Close()
				},
				"server": func(ctx context.Context) error {
					log.Println("Shutting down HTTP server...")
					return httpServer.Shutdown(ctx)
				},
			},
		)
		<-wait
		// End Graceful Shutdown //
	}()

	log.Printf("Server started on port %d", configs.Env.AppPort)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancel()
	wg.Wait()
	log.Println("Server exited")
}

type clientDependency struct {
	// If Needed, add client dependencies here
}

func initClientDependency() client.IRegistry {
	return &clientDependency{
		// Initialize client dependencies here
	}	
}

type repositoryDependency struct {
	DB *sqlx.DB
}

func initRepositoryDependency(dependency repositoryDependency) repository.IRegistry {
	masterTx := repository.NewTransactionRunner(dependency.DB)
	
	repositoryRegistry := repository.NewRegistry(
		masterTx,
	)

	return repositoryRegistry
}

type serviceDependency struct {
	repoRegistry repository.IRegistry
	DB 		 *sqlx.DB
}

func initServiceDependency(dependency serviceDependency) service.IRegistry {
	healthService := serviceHealth.NewHealth(dependency.DB)

	serviceRegistry := service.NewRegistry(
		healthService,
	)

	return serviceRegistry
}

type deliveryDependency struct {
	serviceRegistry *service.IRegistry
}

func initDeliveryDependency(dependency deliveryDependency) delivery.IRegistry {
	registryDelivery := delivery.NewRegistry()
	return registryDelivery
}