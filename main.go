package main

import (
	"VincentLimarus/stock-analyzer-performance/client"
	configs "VincentLimarus/stock-analyzer-performance/config"
	"VincentLimarus/stock-analyzer-performance/migrations"
	"VincentLimarus/stock-analyzer-performance/repository"
	"VincentLimarus/stock-analyzer-performance/repository/trades"
	"VincentLimarus/stock-analyzer-performance/service"
	"VincentLimarus/stock-analyzer-performance/util"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cmdHttp "VincentLimarus/stock-analyzer-performance/cmd/http"
	delivery "VincentLimarus/stock-analyzer-performance/delivery/http"
	serviceHealth "VincentLimarus/stock-analyzer-performance/service/health"
	serviceTrade "VincentLimarus/stock-analyzer-performance/service/trades"

	deliveryTrade "VincentLimarus/stock-analyzer-performance/delivery/http/trades"

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

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		configs.Env.DBHost, configs.Env.DBUser, configs.Env.DBPassword, configs.Env.DBName, configs.Env.DBPort)

	var err error

	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer DB.Close() // ensure DB closes on exit

	log.Println("Connected to PostgreSQL")

	// Auto Migration
	if configs.Env.AutoMigration == "true" {
		if err := migrations.Run(context.Background()); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
	}

	// Initialize Dependencies
	repositoryRegistry := initRepositoryDependency(repositoryDependency{
		DB: DB,
	})

	serviceRegistry := initServiceDependency(serviceDependency{
		repoRegistry: repositoryRegistry,
		DB:           DB,
	})

	deliveryRegistry := initDeliveryDependency(deliveryDependency{
		serviceRegistry: serviceRegistry,
	})

	middlewareLimiter := util.CreateRateLimiter(100, time.Minute)

	// Start HTTP Server
	httpServer := cmdHttp.NewServer(
		deliveryRegistry,
		middlewareLimiter,
	)

	// PERBAIKAN: Jalankan server di goroutine terpisah
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Serve(ctx)
	}()

	log.Printf("Server started on port %s", configs.Env.AppPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancel() // Cancel context untuk signal shutdown

	shutdownDelay := 10 * time.Second
	if configs.Env.ShutDownDelayInSeconds > 0 {
		shutdownDelay = time.Duration(configs.Env.ShutDownDelayInSeconds) * time.Second
	}

	wait := util.GracefullyShutdown(ctx, shutdownDelay,
		map[string]util.Operation{
			"server": func(ctx context.Context) error {
				log.Println("Shutting down HTTP server...")
				return httpServer.Shutdown(ctx)
			},
			"db": func(_ context.Context) error {
				log.Println("Closing database connection...")
				return DB.Close()
			},
		},
	)
	<-wait

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
	tradeRepository := trades.NewTrade(dependency.DB)

	repositoryRegistry := repository.NewRegistry(
		masterTx,
		tradeRepository,
	)

	return repositoryRegistry
}

type serviceDependency struct {
	repoRegistry repository.IRegistry
	DB           *sqlx.DB
}

func initServiceDependency(dependency serviceDependency) service.IRegistry {
	healthService := serviceHealth.NewHealth(dependency.DB)
	tradeService := serviceTrade.NewTrade(dependency.repoRegistry)

	serviceRegistry := service.NewRegistry(
		healthService,
		tradeService,
	)

	return serviceRegistry
}

type deliveryDependency struct {
	serviceRegistry service.IRegistry
}

func initDeliveryDependency(dependency deliveryDependency) delivery.IRegistry {
	tradeDelivery := deliveryTrade.NewTrade(dependency.serviceRegistry)

	registryDelivery := delivery.NewRegistry(
		tradeDelivery,
	)
	return registryDelivery
}
