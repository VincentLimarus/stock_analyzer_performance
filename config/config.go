package configs

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	Env EnvFlat
	once sync.Once
)

type EnvFlat struct {
	AutoMigration string `env:"AUTO_MIGRATION"`
	AppPort	  string `env:"APP_PORT"`
	AppReadHeaderTimeoutInSeconds int `env:"AppReadHeaderTimeoutInSeconds"`
	
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`

	ShutDownDelayInSeconds int `env:"ShutDownDelayInSeconds"`
}

func Init() {
	once.Do(func() {
		// Load .env file
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}

		// Map environment variables to struct
		Env = EnvFlat{
			AutoMigration: os.Getenv("AUTO_MIGRATION"),
			DBHost:        os.Getenv("DB_HOST"),
			DBPort:        os.Getenv("DB_PORT"),
			DBUser:        os.Getenv("DB_USER"),
			DBPassword:    os.Getenv("DB_PASSWORD"),
			DBName:        os.Getenv("DB_NAME"),
			AppPort:       os.Getenv("APP_PORT"),
			AppReadHeaderTimeoutInSeconds: func() int {
				val := os.Getenv("AppReadHeaderTimeoutInSeconds")
				if val == "" {
					return 10 
				}
				var intVal int
				fmt.Sscanf(val, "%d", &intVal)
				return intVal
			}(),
			ShutDownDelayInSeconds: func() int {
				val := os.Getenv("ShutDownDelayInSeconds")
				if val == "" {
					return 15
				}
				var intVal int
				fmt.Sscanf(val, "%d", &intVal)
				return intVal
			}(),
		}

		log.Println("Environment variables loaded successfully")
	})
}
